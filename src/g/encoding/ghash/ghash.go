// 封装常用的hash函数
package ghash


// BKDR Hash Function
func BKDRHash(str []byte) uint32 {
    var seed uint32 = 131; // 31 131 1313 13131 131313 etc..
    var hash uint32 = 0;
    for i := 0; i < len(str); i++ {
        hash = hash * seed + uint32(str[i])
    }
    return (hash & 0x7FFFFFFF)
}

// BKDR Hash Function 64
func BKDRHash64(str []byte) uint64 {
    var seed uint64 = 131; // 31 131 1313 13131 131313 etc..
    var hash uint64 = 0;
    for i := 0; i < len(str); i++ {
        hash = hash * seed + uint64(str[i])
    }
    return (hash & 0x7FFFFFFFFFFFFFFF)
}

// SDBM Hash
func SDBMHash(str []byte) uint32 {
    var hash uint32 = 0;
    for i := 0; i < len(str); i++ {
        // equivalent to: hash = 65599*hash + uint32(str[i]);
        hash = uint32(str[i]) + (hash << 6) + (hash << 16) - hash;
    }
    return (hash & 0x7FFFFFFF);
}

// SDBM Hash 64
func SDBMHash64(str []byte) uint64 {
    var hash uint64 = 0;
    for i := 0; i < len(str); i++ {
        // equivalent to: hash = 65599*hash + uint32(str[i]);
        hash = uint64(str[i]) + (hash << 6) + (hash << 16) - hash;
    }
    return (hash & 0x7FFFFFFFFFFFFFFF);
}

// RS Hash Function
func RSHash(str []byte) uint32 {
    var b uint32 = 378551;
    var a uint32 = 63689;
    var hash uint32 = 0;
    for i := 0; i < len(str); i++ {
        hash = hash * a + uint32(str[i]);
        a *= b;
    }
    return (hash & 0x7FFFFFFF);
}

// RS Hash Function 64
func RSHash64(str []byte) uint64 {
    var b    uint64 = 378551;
    var a    uint64 = 63689;
    var hash uint64 = 0;
    for i := 0; i < len(str); i++ {
        hash = hash * a + uint64(str[i]);
        a *= b;
    }
    return (hash & 0x7FFFFFFFFFFFFFFF);
}

// JS Hash Function
func JSHash(str []byte) uint32 {
    var hash uint32 = 1315423911;
    for i := 0; i < len(str); i++ {
        hash ^= ((hash << 5) + uint32(str[i]) + (hash >> 2));
    }
    return (hash & 0x7FFFFFFF);
}

// JS Hash Function 64
func JSHash64(str []byte) uint64 {
    var hash uint64 = 1315423911;
    for i := 0; i < len(str); i++ {
        hash ^= ((hash << 5) + uint64(str[i]) + (hash >> 2));
    }
    return (hash & 0x7FFFFFFFFFFFFFFF);
}

// P. J. Weinberger Hash Function
func PJWHash(str []byte) uint32 {
    var BitsInUnignedInt uint32 = (4 * 8);
    var ThreeQuarters    uint32 = ((BitsInUnignedInt  * 3) / 4);
    var OneEighth        uint32 = (BitsInUnignedInt / 8);
    var HighBits         uint32 = (0xFFFFFFFF) << (BitsInUnignedInt - OneEighth);
    var hash uint32             = 0;
    var test uint32             = 0;
    for i := 0; i < len(str); i++ {
        hash = (hash << OneEighth) + uint32(str[i]);
        if test = hash & HighBits; test != 0 {
            hash = ((hash ^ (test >> ThreeQuarters)) & (^HighBits + 1));
        }
    }
    return (hash & 0x7FFFFFFF);
}

// P. J. Weinberger Hash Function 64
func PJWHash64(str []byte) uint64 {
    var BitsInUnignedInt uint64 = (4 * 8);
    var ThreeQuarters    uint64 = ((BitsInUnignedInt  * 3) / 4);
    var OneEighth        uint64 = (BitsInUnignedInt / 8);
    var HighBits         uint64 = (0xFFFFFFFF) << (BitsInUnignedInt - OneEighth);
    var hash             uint64 = 0;
    var test             uint64 = 0;
    for i := 0; i < len(str); i++ {
        hash = (hash << OneEighth) + uint64(str[i]);
        if test = hash & HighBits; test != 0 {
            hash = ((hash ^ (test >> ThreeQuarters)) & (^HighBits + 1));
        }
    }
    return (hash & 0x7FFFFFFFFFFFFFFF);
}

// ELF Hash Function
func ELFHash(str []byte) uint32 {
    var hash uint32 = 0;
    var x    uint32 = 0;
    for i := 0; i < len(str); i++ {
        hash = (hash << 4) + uint32(str[i]);
        if x = hash & 0xF0000000; x != 0 {
            hash ^= (x >> 24);
            hash &= ^x + 1;
        }
    }
    return (hash & 0x7FFFFFFF);
}

// ELF Hash Function 64
func ELFHash64(str []byte) uint64 {
    var hash uint64 = 0;
    var x    uint64 = 0;
    for i := 0; i < len(str); i++ {
        hash = (hash << 4) + uint64(str[i]);
        if x = hash & 0xF0000000; x != 0 {
            hash ^= (x >> 24);
            hash &= ^x + 1;
        }
    }
    return (hash & 0x7FFFFFFFFFFFFFFF);
}

// DJB Hash Function
func DJBHash(str []byte) uint32 {
    var hash uint32 = 5381;
    for i := 0; i < len(str); i++ {
        hash += (hash << 5) + uint32(str[i]);
    }
    return (hash & 0x7FFFFFFF);
}

// DJB Hash Function 64
func DJBHash64(str []byte) uint64 {
    var hash uint64 = 5381;
    for i := 0; i < len(str); i++ {
        hash += (hash << 5) + uint64(str[i]);
    }
    return (hash & 0x7FFFFFFFFFFFFFFF);
}

// AP Hash Function
func APHash(str []byte) uint32 {
    var hash uint32 = 0;
    for i := 0; i < len(str); i++ {
        if ((i & 1) == 0) {
            hash ^= ((hash << 7) ^ uint32(str[i]) ^ (hash >> 3));
        } else {
            hash ^= (^((hash << 11) ^ uint32(str[i]) ^ (hash >> 5)) + 1);
        }
    }
    return (hash & 0x7FFFFFFF);
}

// AP Hash Function 64
func APHash64(str []byte) uint64 {
    var hash uint64 = 0;
    for i := 0; i < len(str); i++ {
        if ((i & 1) == 0) {
            hash ^= ((hash << 7) ^ uint64(str[i]) ^ (hash >> 3));
        } else {
            hash ^= (^((hash << 11) ^ uint64(str[i]) ^ (hash >> 5)) + 1);
        }
    }
    return (hash & 0x7FFFFFFFFFFFFFFF);
}