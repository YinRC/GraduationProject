package p_machine

import (
    "fmt"
)

// 遍历切片的每个元素, 通过给定函数进行元素访问
func visit(list []int, f func(int)) {

    for _, v := range list {
        f(v)
    }
}
