import React, {useState} from "react";


const FileUploader = () => {

    const sendFile = () => {
        const data = new FormData();
        data.set("file", document.getElementById("file").files[0]);
        fetch("/api/test", {
            method: "POST",
            body: data
        }).then(response => response.json())
            .then(response => console.log(response));
    }

    return (
        <div>
            <h2>Upload a file</h2>
            <label for="file">Filename:</label>
            <input type="file" name="file" id="file"/>
            <br/>
            <input type="submit" name="submit" value="Submit" onClick={sendFile}/>

        </div>
    );
};

export default FileUploader;
