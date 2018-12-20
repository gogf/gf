package main

type Registry struct {
    Method  string
    Uri     string
    Handler interface{}
    Object  interface{}
}

func BindGroup(group string, routers [][]interface{}) {

}

type User    struct { }
type Order   struct { }
type Product struct { }

func HookFunc() {

}

func main() {
    user := new(User)
    BindGroup("/api", [][]interface{} {
        {"ALL",   "/*",               HookFunc, "BeforeServe"},
        {"ALL",   "/order",           new(Order)},
        {"REST",  "/product",         new(Product)},
        {"GET",   "/user/register",   "Register",      user},
        {"GET",   "/user/reset-pass", "ResetPassword", user},
        {"POST",  "/user/reset-pass", "ResetPassword", user},
        {"POST",  "/user/login",      "Login",         user},
    })
}
