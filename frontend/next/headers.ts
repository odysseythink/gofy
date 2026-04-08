export async function cookies() {
  return {
    get: (name: string) => {
      const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]*)`))
      return match ? { name, value: decodeURIComponent(match[1]) } : undefined
    },
    set: () => {},
    delete: () => {},
  }
}

export async function headers() {
  const map = new Map<string, string>()
  return {
    get: (key: string) => map.get(key) ?? null,
    forEach: (_cb: (value: string, key: string) => void) => {},
  }
}
