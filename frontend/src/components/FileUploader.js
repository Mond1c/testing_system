import React, { useEffect, useState } from "react";

const FileUploader = () => {
  const [problems, setProblems] = useState([]);
  const [username, setUsername] = useState(undefined);
  const [languages, setLanguages] = useState([]);

  const getProblems = () => {
    fetch("/api/problems")
      .then((response) => response.json())
      .then((response) => {
        response.problems.sort();
        setProblems(response.problems);
      });
  };

  const getLanguages = () => {
    fetch("/api/languages")
      .then(response => response.json())
      .then(response => setLanguages(response));
  };

  const getUsername = () => {
    fetch("/api/me")
      .then((response) => response.json())
      .then((response) => setUsername(response.username));
  };

  const sendFile = () => {
    if (username === undefined) {
      console.error("Username is undefined");
      return;
    }
    const data = new FormData();
    data.set("file", document.getElementById("file").files[0]);
    data.set("language", document.getElementById("language").value);
    data.set("problem", document.getElementById("problem").value);
    data.set("username", username);
    fetch("/api/test", {
      method: "POST",
      body: data,
    })
      .then((response) => response.json())
      .then((response) => console.log(response));
  };

  useEffect(() => {
    getUsername();
    getProblems();
    getLanguages();
  }, []);

  return (
    <div>
      <h1>Hello, {username !== undefined ? username : ""}</h1>
      <h2>Upload a solution</h2>
      <label for="problem">Problem:</label>
      <select name="problem" id="problem">
        {problems.map((problem) => {
          return <option value={problem}>{problem}</option>;
        })}
      </select>
      <br />
      <label for="file">Filename:</label>
      <input type="file" name="file" id="file" />
      <br />
      <input type="submit" name="submit" value="Submit" onClick={sendFile} />
      <br />
      <label for="language">Language:</label>
      <select name="language" id="language">
        {languages.map(lang => {
          return <option value={lang.value}>{lang.name}</option>;
        })}
      </select>
    </div>
  );
};

export default FileUploader;
