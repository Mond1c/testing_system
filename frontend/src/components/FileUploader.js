import React from "react";


const FileUploader = () => {

    const sendFile = () => {
        const data = new FormData();
        data.set("file", document.getElementById("file").files[0]);
        data.set("language", document.getElementById("language").value);
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
            <select name="language" id="language">
                <option value="cpp">C++ 20</option>
                <option value="java">Java 21</option>
                <option value="go">Go 1.21</option>
            </select>
        </div>
    );
};

export default FileUploader;
