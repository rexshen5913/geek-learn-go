package web

import (
	"fmt"
	"regexp"
	"strings"
)

// 用來支持對路由樹的操作
// 代表路由樹的(森林)
type router struct {
	// Beego Gin HTTP method 對應一棵樹
	// GET 有一棵樹, POST 有一棵樹

	// http method -> 路由樹根節點
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: make(map[string]*node),
	}
}

// 這邊需要加上校驗，避免用戶輸入不適當的 path
// 像是 ///user//a//b/c/d///
// path 必須以 / 開頭，不能以 / 結尾，不能中間有連續的 //
// 不能在同一個位置同時註冊 * 和 : 路由
// 不能在同一個位置同時註冊多個參數路由
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	// 校驗 path
	if path == "" {
		panic("web: 路由不能為空字符串")
	}
	if path[0] != '/' {
		panic("web: 路由必須以 / 開頭")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 結尾")
	}
	// 這邊需要校驗連續的 //
	// 可以試著用 regexp 來校驗
	re := regexp.MustCompile(`//+`)
	if re.MatchString(path) {
		panic("web: 路由不能包含連續多個 //")
	}

	// 首先找到樹來
	root, ok := r.trees[method]
	if !ok {
		// 還沒有根節點，先創建一個根節點
		root = &node{
			path:     "/",
			children: make(map[string]*node),
		}
		r.trees[method] = root
	}

	// 跟節點特殊處理下
	// 跟節點重複註冊
	if path == "/" {
		if root.handler != nil {
			panic("web: 路由不能重複註冊 [/]")
		}
		root.handler = handleFunc
		root.route = "/"
		return
	}

	// root.addRoute(path, handleFunc)

	// 切割這個 path
	// 如果 path 是 /user/home ，那麼 segs 就是 ["user", "home"]
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")

	// 遍歷 segs
	for _, seg := range segs {
		// 遞歸下去，找准位置
		// 如果中途有節點不存在，你就要創建出來
		children := root.childOfCreate(seg)
		root = children
	}
	if root.handler != nil {
		panic("web: 路由不能重複註冊 [" + path + "]")
	}
	root.handler = handleFunc
	root.route = path
	fmt.Printf("注册路由后的树结构:\n")
	r.printTree()
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	fmt.Printf("查找路由: %s %s\n", method, path)
	// 基本上也是沿著樹深度查找下去
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}

	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	var pathParams map[string]string

	for i, seg := range segs {
		fmt.Printf("当前节点: %s, 查找段: %s\n", root.path, seg)
		child, paramChild, found := root.childOf(seg)
		if !found {
			fmt.Printf("未找到匹配，检查通配符节点\n")
			// 如果沒有找到匹配的子節點，檢查是否有通配符節點
			if root.starChild != nil {
				fmt.Printf("找到通配符节点，返回\n")
				// 找到通配符節點，直接返回，因為通配符可以匹配剩餘的所有路徑段
				return &matchInfo{
					n:          root.starChild,
					pathParams: pathParams,
				}, true
			}
			return nil, false
		}
		if paramChild {
			// 命中了參數路由，需要把參數路由的值記錄下來
			// 這邊需要記錄下來，後續用戶請求的時候，需要用到
			// path 是 :id 的形式
			if pathParams == nil {
				pathParams = make(map[string]string)
			}

			// 根据节点类型设置参数名
			if child.regex != nil {
				pathParams[child.regexParam] = seg
			} else {
				pathParams[child.path[1:]] = seg
			}
		}
		root = child

		// 这个逻辑是必要的！
		if root.path == "*" && i < len(segs)-1 {
			fmt.Printf("当前是通配符节点，还有剩余路径段，直接返回\n")
			return &matchInfo{
				n:          root,
				pathParams: pathParams,
			}, true
		}
	}
	// 這樣return是說確實有這節點，
	// 但是無法保證這邊的節點是用戶註冊且有 handler 的
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
}

func (n *node) childOfCreate(seg string) *node {

	// 处理通配符
	if seg == "*" {
		if n.paramChild != nil || n.regexChild != nil {
			panic("web: 路由不能同时注册 * 和参数路由")
		}
		if n.starChild != nil {
			return n.starChild
		}
		n.starChild = &node{path: seg}
		return n.starChild
	}

	// 处理参数路由（包括正则表达式）
	if paramName, regexPattern, isParam := parseParamRoute(seg); isParam {
		if n.starChild != nil {
			panic("web: 路由不能同时注册 * 和参数路由")
		}

		if regexPattern != "" {
			// 正则表达式参数路由
			if n.paramChild != nil || n.regexChild != nil {
				panic("web: 路由不能同时注册多个参数路由")
			}

			// 编译正则表达式
			regexPattern = "^" + regexPattern + "$"
			regex, err := regexp.Compile(regexPattern)
			if err != nil {
				panic(fmt.Sprintf("web: 无效的正则表达式 [%s]: %v", regexPattern, err))
			}

			n.regexChild = &node{
				path:       seg,
				regexParam: paramName,
				regex:      regex,
			}
			return n.regexChild
		} else {
			// 普通参数路由
			if n.paramChild != nil || n.regexChild != nil {
				panic("web: 路由不能同时注册多个参数路由")
			}
			n.paramChild = &node{path: seg}
			return n.paramChild
		}
	}

	// 处理静态路由
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	child, ok := n.children[seg]
	if !ok {
		child = &node{path: seg}
		n.children[seg] = child
	}
	return child
}

type node struct {
	route string

	path string

	// 静态节点
	children map[string]*node

	// 业务逻辑
	handler HandleFunc

	// 通配符 * 匹配的节点
	starChild *node

	// 参数路由
	paramChild *node

	// 正则表达式参数路由
	regexChild *node
	regex      *regexp.Regexp
	regexParam string
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

// 解析参数路由，支持正则表达式
// 例如：:id(\d+) -> paramName: "id", regex: "\d+"
func parseParamRoute(seg string) (string, string, bool) {
	if len(seg) < 2 || seg[0] != ':' {
		return "", "", false
	}

	// 支援 :paramName 或 :paramName(regex)
	re := regexp.MustCompile(`^:([a-zA-Z_][a-zA-Z0-9_]*)(?:\(([^()]+)\))?$`)
	matches := re.FindStringSubmatch(seg)

	if len(matches) == 0 {
		return "", "", false
	}

	param := matches[1]
	regex := matches[2] // 若沒有正則，這裡會是空字串
	return param, regex, true

}

// childOf 優先考慮靜態匹配，如果沒有靜態匹配，則考慮通配符匹配
// 返回值：	子節點，是否是參數路由，命中了沒有
func (n *node) childOf(seg string) (*node, bool, bool) {

	// 优先静态匹配
	if n.children != nil {
		if child, ok := n.children[seg]; ok {
			return child, false, true
		}
	}

	// 正则表达式参数匹配
	if n.regexChild != nil && n.regexChild.regex != nil && n.regexChild.regex.MatchString(seg) {
		return n.regexChild, true, true
	}

	// 普通参数匹配
	if n.paramChild != nil {
		return n.paramChild, true, true
	}

	// 通配符匹配
	return n.starChild, false, n.starChild != nil
}
func (r *router) printTree() {
	for method, root := range r.trees {
		fmt.Printf("Method: %s\n", method)
		printNode(root, 0)
	}
}

func printNode(n *node, depth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s%s (handler: %v)\n", indent, n.path, n.handler != nil)
	for _, child := range n.children {
		printNode(child, depth+1)
	}
	if n.starChild != nil {
		printNode(n.starChild, depth+1)
	}
	if n.paramChild != nil {
		printNode(n.paramChild, depth+1)
	}
}
