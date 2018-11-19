# safemap
safemap is a simple thread-safe map in golang
sample:

```
	m := util.NewSmap()

	for i := 0; i < 5; i++ {
		m.Set(i, i + 1)
	}
	fmt.Println(m.Size())
	val, ok := m.Get(1)
	if ok {
		fmt.Println(val)
	}
	m.Del(1)
	fmt.Println(m.Size())
```
