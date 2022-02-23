package utils

import (
	"bytes"
	"strings"
)

const (
	static nodeType = iota // default
	root
	param
	query
)

type nodeType uint8

type node struct {
	// 节点路径，比如上面的s，earch，和upport
	path string
	// 儿子节点
	children []*node
	nType    nodeType
	// 完整路径
	tokenAuth       string // token 校验
	permissionsAuth string // 权限校验
}

type methodTree struct {
	method string
	root   *node
}

type Engine struct {
	trees methodTrees
}

type methodTrees []methodTree

// 通过mothod 去获取该mothod 的路由树
func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

// 构建路由树
func (n *node) addRoute(path, tokenAuth, permissionsAuth string) {
	pathList := strings.Split(path, "/") // 获取url
	s := []byte(path)
	countSlash := uint16(bytes.Count(s, []byte("/"))) // 获取'/'的个数
	if countSlash == 0 {
		return
	}
	// 如果路由中只有一个'/',那么就直接赋值给根路径
	if countSlash == 1 && len(pathList) == 0 {
		n.nType = root
		n.path = "/"
		n.tokenAuth = tokenAuth
		n.permissionsAuth = permissionsAuth
	} else {
		// 构建子路由树
		n.insertChild(path, tokenAuth, permissionsAuth)
	}

}

// 获取url的列表,将位置参数和关键字参数都添加到列表中
func getUrlList(path string) []string {
	pathList := strings.Split(path, "/") // 获取url
	list := []string{}
	for _, p := range pathList { // 重新构造一个列表,其中包含位置参数和关键字参数
		if p == "" {
			continue
		}
		// 判断路径里面是否存在 关键字参数
		index := bytes.IndexByte([]byte(p), '?')
		if index == -1 {
			list = append(list, p)
		} else {
			// 将关键字参数作为一层子树放进列表中
			list = append(list, p[:index], p[index:])
		}
	}
	return list
}

func (n *node) insertChild(path, tokenAuth, permissionsAuth string) {
	list := getUrlList(path)
	head := n // 指向头节点的children

	llen := len(list)
	// 遍历路由列表和每层路由树
	for index1, l := range list {
		findflag := false
		// 开始遍历子树,从跟路由的子树开始遍历
		for index2, n1 := range head.children {
			// 当前子树中存在该路由
			if n1.path == l {
				if llen == index1+1 { // 遍历到了最后一个路径,将tokenAuth和permissionsAuth 赋值给该节点的tokenAuth和permissionsAuth
					n1.tokenAuth = tokenAuth
					n1.permissionsAuth = permissionsAuth
				}
				head = head.children[index2] // 指向当前节点的下一个children,继续遍历
				findflag = true
				break
			}
		}
		if findflag {
			// 表示子树中存在路由,且路由还没有匹配完
			continue
		}
		// 该层子树中没有该路径,那么就添加该路径到该层子树下面,并且以新建的子树创建下层子树
		var nType nodeType
		// 识别当前是什么类型的节点
		if l[0] != ':' && l[0] != '?' {
			nType = static
		} else if l[0] == ':' {
			nType = param
		} else if l[0] == '?' {
			nType = query
			l = l[1:] // 去掉字符串中的?
		}
		if llen == index1+1 { // 遍历到了最后一个路径,将tokenAuth和permissionsAuth 赋值给该节点的tokenAuth和permissionsAuth
			head.children = append(head.children, &node{path: l, nType: nType, tokenAuth: tokenAuth, permissionsAuth: permissionsAuth})
		} else {
			head.children = append(head.children, &node{path: l, nType: nType})
		}
		hlen := len(head.children)
		head = head.children[hlen-1] // 指向当前节点的下一个children,创建下层子树
	}
}

// 构建路由树
func (engine *Engine) addRoute(method, path, tokenAuth, permissionsAuth string) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")

	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.path = "/"
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, tokenAuth, permissionsAuth)
}

func (engine *Engine) post(relativePath, tokenHandle, persissionHandle string) {
	engine.addRoute("post", relativePath, tokenHandle, persissionHandle)
}

func (engine *Engine) get(relativePath, tokenHandle, persissionHandle string) {
	engine.addRoute("get", relativePath, tokenHandle, persissionHandle)
}

func (engine *Engine) delete(relativePath, tokenHandle, persissionHandle string) {
	engine.addRoute("delete", relativePath, tokenHandle, persissionHandle)
}

func (engine *Engine) patch(relativePath, tokenHandle, persissionHandle string) {
	engine.addRoute("patch", relativePath, tokenHandle, persissionHandle)
}

func (engine *Engine) update(relativePath, tokenHandle, persissionHandle string) {
	engine.addRoute("update", relativePath, tokenHandle, persissionHandle)
}

/*
解析路由树
url: 解析的路径, method: 请求方法, querys: 关键字参数
return  tokenauth,验证身份的函数,为空则表示不需要验证
*/
func (engine *Engine) ParseUrlTree(url, method, querys string) string {
	// 是否存在该方法的路由树
	root := engine.trees.get(method)
	if root == nil {
		return ""
	}
	if url == "/" {
		if root.path == url {
			return root.tokenAuth
		}
		return ""
	}
	// 构建路由的队列
	list := getUrlList(url)
	if querys != "" {
		list = append(list, querys)
	}
	llen := len(list)
	for index1, p := range list {
		compareFlag := false
		index3 := 0
		// 从根路径的子树开始遍历匹配
		for index2, p1 := range root.children {
			index3++
			// 如果类型是关键字参数类型,则比较路由树里面的参数是否与路由列表的参数匹配
			if p1.nType == query {
				qflag := true
				for _, q := range joinUrlQuery(p1.path) {
					if !strings.Contains(p, q) {
						qflag = false
						break
					}
				}
				if qflag { // 如果query参数全部包含则表示匹配成功
					compareFlag = true
				}
			} else {
				if p1.path == p {
					compareFlag = true
				}
			}
			// 匹配到路由了判断是否是最后一个字段,是的话,返回验证token 的函数字符串
			if compareFlag {
				if llen == index1+1 {
					return p1.tokenAuth
				}
				root = root.children[index2]
				break
			}
		}
		if index3 == len(root.children) && !compareFlag { // 当前child全部遍历都没有找到相匹配的节点就直接退出去
			break
		}
	}
	return ""
}

func joinUrlQuery(urlQuery string) []string {
	byteUrlQuery := []byte(urlQuery)
	flag := false
	byteQuey := []byte{}
	for _, b := range byteUrlQuery {
		if b == '=' {
			flag = true
		} else if b == '&' {
			flag = false
		}
		if flag {
			continue
		}
		byteQuey = append(byteQuey, b)
	}
	return strings.Split(string(byteQuey), "&")
}

var Root Engine

// 初始化路由树
func InitTree() {
	// 参数1表示路由,路由规则:如果存在位置参数(param),那么与注册路由时写法保持一致;如果存在query,那么接'?'和query名,有多个则用'&'拼接.
	// 参数2表示的是需要使用哪种身份验证函数,会自动根据你填写的字符串去匹配验证方法
	// 参数3目前还没有写,本意是打算做权限验证
	Root.get("/test/:id", "tokenAuth1", "")
	Root.get("/test/t1/:id", "tokenAuth1", "")
	Root.get("/test/t2", "", "")
	Root.get("/test/t/:id1?name&pwd", "tokenAuth1", "")
	Root.get("/test/t2/:id1/:id2", "", "")
	Root.get("/test", "", "")

	Root.post("/test/:id", "", "")
	Root.post("/test/t1/:id", "", "")
	Root.post("/test/t2", "", "")
	Root.post("/test/t/:id1", "", "")
	Root.post("/test/t2/:id1/:id2", "", "")
	Root.post("/test", "", "")
}
