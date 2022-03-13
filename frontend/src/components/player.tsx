import React from "react";

import {getExtension} from "../utils";


const playerTypeVideo: { [name: string]: string } = {
    "mkv": "mp4",
    "avi": "mp4"
}

interface PlayerProps {
    file: string | null
    setState?: any
}

const PlayerView: React.FC<any> = (props: PlayerProps) => {
    const {file, setState} = props
    let type = getExtension(file)
    type = playerTypeVideo[type] ? playerTypeVideo[type] : type
    const url = `/video/?path=${file}`

    return <video width="100%" height="100%" controls autoPlay style={{padding: 0, margin: 0}} onPlay={setState}>
        <source src={url} type={`video/${type}`}/>
        Your browser does not support the video tag.
    </video>
}

export default PlayerView