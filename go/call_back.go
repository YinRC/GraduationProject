package p_machine

// 遍历切片的每个元素, 通过给定函数进行元素访问
func Visit(list []int, f func(int)) {

    for _, v := range list {
        f(v)
    }
}
