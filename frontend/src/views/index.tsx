import React, {useMemo} from "react";
import {getExtension} from "../utils";

import axios from "axios";
import PlayerView from "../components/player";

const ITEM_LENGTH = 30;

const Extension_dict: { [name: string]: string } = {
    "folder": "folder.png",
    "mp4": "premium-icon-mp4-file-2266843.png",
    "avi": "premium-icon-avi-file-2266783.png",
    "mkv": "mkv.png"
}

interface ItemProps {
    name: string
    isDir: boolean
}

const getImageExtension = (filename: string) => {
    const extension = getExtension(filename)
    const image: string = Extension_dict[extension]
    if (image) return image
    return "file.png"
}

function Index() {
    const {useState, useEffect} = React
    const [files, setFiles] = useState([])
    const [file, setFile] = useState<string>("")
    const [currentPath, setCurrentPath] = useState("")
    const [loading, setLoading] = useState(false)

    const getFiles = (path = "") => {
        const url =`/api/list/?path=${encodeURI(path)}`
        axios.get(url)
            .then(r => {
                if (r.status === 200) {
                    const json = r.data
                    return json
                }
            })
            .then(setFiles)
    }

    const loadClip = (file: string) => {
        setLoading(true)
        setFile(file)
        setTimeout(() => {
            setLoading(false)
        }, 1000)
    }

    const openPath = (path: string) => {
        if (path === "..") {
            const newPathSplit = currentPath.split("/")
            path = newPathSplit.slice(0, newPathSplit.length - 1).join("/")
        } else {
            path = currentPath + "/" + path
        }
        setCurrentPath(path)
        getFiles(path)
    }

    useEffect(() => {
        getFiles()
    }, [])

    const LeftSideBlock: React.FC<any> = ({children}: any) => {
        return <div style={{
            width: "300px",
            height: "100vh",
            padding: "0",
            borderRight: "solid 1px white",
            backgroundColor: "black"
        }}>
            <div style={{padding: "15px 10px"}}>
                {children}
            </div>
        </div>
    }

    const RightSideBlock: React.FC<any> = ({children}: any) => {
        return useMemo(() => <div style={{flex: 1, backgroundColor: "#000", height: "100vh", width: "100%"}}>
            {children}
        </div>, [children])
    }

    const ItemFile = ({path, onClick}: { path: ItemProps, onClick: any }) => {
        const filename = path.name
        const extension = path.isDir ? "folder.png" : getImageExtension(filename)
        const isSelected = file.indexOf(filename)
        const style = isSelected > 0 ? {backgroundColor: "blue"} : {}
        console.log(file, filename)

        return <div style={{padding: "4px", ...style}}>
            <a style={{color: "whitesmoke", textDecoration: "none"}} onClick={onClick} title={filename} href="#">
                <div style={{display: "flex", flexDirection: "row"}}>
                    <img width="20px" height="20px" src={`/static/images/${extension}`} style={{marginRight: "5px"}}/>
                    {filename.substring(0, ITEM_LENGTH)}{filename.length > ITEM_LENGTH ? "..." : ""}
                </div>
            </a>
        </div>
    }

    const FilesBlock: React.FC<any> = (props: { files: ItemProps[] }) => {
        return <div>
            {props.files?.sort((a, b) => Number(b.isDir) - Number(a.isDir)).map((v, i) =>
                <ItemFile key={i} path={v}
                          onClick={() => v.isDir ? openPath(`${v.name}`) : loadClip(`${currentPath}/${v.name}`)}/>
            )}
        </div>
    }

    return (
        <div style={{display: "flex", flexDirection: "row", width: "100%", height: "100%", padding: 0, margin: 0}}>
            <LeftSideBlock>
                <FilesBlock files={files}/>
            </LeftSideBlock>
            <RightSideBlock>
                {!loading && file && <PlayerView file={file}/>}
            </RightSideBlock>
        </div>
    )
}

export default Index