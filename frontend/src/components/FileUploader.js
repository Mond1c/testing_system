import React, {useState} from "react";


const FileUploader = () => {
    return (
        <div>
            <h2>Upload a file</h2>
            <form action="/api/test" method="post" enctype="multipart/form-data">
                <label for="file">Filename:</label>
                <input type="file" name="file" id="file"/>
                <br/>
                <input type="submit" name="submit" value="Submit"/>
            </form>
        </div>
    );
};

export default FileUploader;
