package packed

import "github.com/gogf/gf/os/gres"

func init() {
	if err := gres.Add("H4sIAAAAAAAC/wrwZmYRYeBg4GBIi7kQwIAE+Bk4GZLz89Iy0/VLi3L0SvJzc0JDWBkYC9a9iMvpc+xrNuBxvf5DJH7m2xM3+494KolpqLRpHtFgWbm3dmpI1u6E5+//z3sq73HE7pf0nmU7PGYfTonvdO1cYbg5j/F4ENOfKTPmzL/kwFDwTNT2TMypuj+R2hu//bdST5zvs9X+1cxwtfDPW1cb9kcnnL3kEvgm8j2/y9mv7JytJ7devW1/4XDg9VrZTAb9rLPvQx9ONr3Ce56BgeH//wBvdo4D12YyTmNgYPjFwMAA8xsDw+EMwUBkv7HB/Qb20v4Gq3iQZmQlAd6MTCLMiKBBNhgUNDCwpBFE4goohCnYHQEBAgz/He/ATUFyEisbSJqJgYmhmYGBQZIRxAMEAAD//5uozoWyAQAA"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
