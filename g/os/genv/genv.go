package genv

import "os"

func All() []string {
    return os.Environ()
}

func Get(k string) string {
    return os.Getenv(k)
}

func Set(k, v string) error {
    return os.Setenv(k, v)
}

func Remove(k string) error {
    return os.Unsetenv(k)
}