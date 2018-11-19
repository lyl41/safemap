# safemap
safemap is a simple thread-safe map in golang
sample:

```
m := util.NewSmap()

	m.Set(1, 2)
	val, ok := m.Get(1)
	fmt.Println(val, ok)
	fmt.Println(m.Size())
```
