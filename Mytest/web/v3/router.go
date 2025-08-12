package web

import (
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
}

func (r *router) findRoute(method string, path string) (*node, bool) {
	// 基本上也是沿著樹深度查找下去
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return root, true
	}

	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")

	for _, seg := range segs {
		child, ok := root.childOf(seg)
		if !ok {
			return nil, false
		}
		root = child
	}
	// 這樣return是說確實有這節點，
	// 但是無法保證這邊的節點是用戶註冊且有 handler 的
	return root, true
}

func (n *node) childOfCreate(seg string) *node {
	// 加上判斷 *
	if seg == "*" {
		if n.starChild != nil {
			// 重複註冊就當沒發生就好
			return n.starChild
		}
		n.starChild = &node{
			path: seg,
		}
		return n.starChild
	}
	// 如果沒有子節點，創建一個子節點
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	child, ok := n.children[seg]
	if !ok {
		child = &node{
			path: seg,
		}
		n.children[seg] = child
	}
	return child
}

type node struct {
	path string

	// 靜態節點
	// 子 path 到子節點的映射
	children map[string]*node

	// 缺一個代表用戶註冊的業務邏輯
	handler HandleFunc

	// 通配符 * 匹配的節點
	starChild *node
}

// childOf 優先考慮靜態匹配，如果沒有靜態匹配，則考慮通配符匹配
func (n *node) childOf(seg string) (*node, bool) {
	if n.children == nil {
		return n.starChild, n.starChild != nil
	}
	child, ok := n.children[seg]
	if !ok {
		return n.starChild, n.starChild != nil
	}
	return child, true
}
