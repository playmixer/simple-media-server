
export const getExtension = (filename) => {
    const type_split = filename.split(".")
    const type = type_split[type_split.length - 1]
    return type
}