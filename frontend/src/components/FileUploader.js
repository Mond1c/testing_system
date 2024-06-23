import React, { useEffect, useState } from "react";
import styled from "styled-components";

const Container = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  margin-left: 20px;
`;

const Title = styled.h1`
  text-align: left;
`;

const Label = styled.label`
  display: block;
  text-align: left;
`;

const Select = styled.select`
  text-align: left;
`;

const Input = styled.input`
  text-align: left;
`;

const Button = styled.input`
  text-align: left;
`;

const FileUploader = () => {
  const [problems, setProblems] = useState([]);
  const [username, setUsername] = useState(undefined);
  const [languages, setLanguages] = useState([]);
  const [verdict, setVerdict] = useState("");

  const getProblems = () => {
    fetch("/api/problems")
      .then((response) => response.json())
      .then((response) => {
        setProblems(response);
      });
  };

  const getLanguages = () => {
    fetch("/api/languages")
      .then((response) => response.json())
      .then((response) => setLanguages(response));
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
      .then((response) => {
        setVerdict(response.message);
      });
  };

  useEffect(() => {
    getUsername();
    getProblems();
    getLanguages();
  }, []);

  return (
    <Container>
      <Title>Hello, {username !== undefined ? username : ""}</Title>
      <h2 style={{ textAlign: "left" }}>Upload a solution</h2>
      <Label htmlFor="problem">Problem:</Label>
      <Select name="problem" id="problem">
        {problems.map((problem) => (
          <option key={problem} value={problem}>
            {problem}
          </option>
        ))}
      </Select>
      <br />
      <Label htmlFor="file">Filename:</Label>
      <Input type="file" name="file" id="file" />
      <br />
      <Label htmlFor="language">Language:</Label>
      <Select name="language" id="language">
        {languages.map((lang) => (
          <option key={lang.value} value={lang.value}>
            {lang.name}
          </option>
        ))}
      </Select>
      <br />
      <Button type="submit" name="submit" value="Submit" onClick={sendFile} />
      <br />
      <Title>Verdict: {verdict}</Title>
    </Container>
  );
};

export default FileUploader;
