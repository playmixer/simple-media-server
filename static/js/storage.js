const store = window.localStorage

function getItem(key) {
    return store.getItem(key)
}

function setItem(k, v) {
    store.setItem(k, v)
}

function getInfo(file) {
    try {
        const info = JSON.parse(getItem(file))
        return info
    } catch {
        return {}
    }
}

function updInfo(file, params) {
    const info = getInfo(file)
    setItem(file, JSON.stringify({
        ...info,
        ...params
    }))
}