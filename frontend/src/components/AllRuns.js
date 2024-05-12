import React, { useEffect, useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { solarizedDark } from "react-syntax-highlighter/dist/esm/styles/hljs";
import { solarizedlight } from "react-syntax-highlighter/dist/esm/styles/prism";
import styled from "styled-components";

const TableWrapper = styled.table`
    border: 1px solid black;
    width: 100%;
`;

const TableHeaderWrapper = styled.th`
    border: 1px solid black;
`;

const TableElementWrapper = styled.td`
    border: 1px solic black;
    text-aligh: center;
`;

const AllRuns = () => {
  const [runs, setRuns] = useState([]);
  const [sourceCode, setSourceCode] = useState(undefined);
  const [language, setLanguage] = useState("");

  const RunInfo = ({ run }) => {
    const viewSourceCode = () => {
      fetch(
        "/api/admin/source_code?username=" + run.username + "&run_id=" +
          run.run_id,
      )
        .then((response) => response.text())
        .then((data) => {
          setSourceCode(data);
          setLanguage(run.language);
        });
    };

    return (
      <tr>
        <TableElementWrapper>{run.username}</TableElementWrapper>
        <TableElementWrapper>{run.problem}</TableElementWrapper>
        <TableElementWrapper>{run.result}</TableElementWrapper>
        <TableElementWrapper>{run.time}</TableElementWrapper>
        <TableElementWrapper>
          <button onClick={viewSourceCode}>View</button>
        </TableElementWrapper>
      </tr>
    );
  };

  const RunTable = ({ runs }) => {
    return (
      <TableWrapper>
        <tr>
          <TableHeaderWrapper>Username</TableHeaderWrapper>
          <TableHeaderWrapper>Problem</TableHeaderWrapper>
          <TableHeaderWrapper>Result</TableHeaderWrapper>
          <TableHeaderWrapper>Time</TableHeaderWrapper>
          <TableHeaderWrapper>Source Code</TableHeaderWrapper>
        </tr>
        {runs.map((run, index) => (
          <RunInfo run={run} key={index} />
        ))}
      </TableWrapper>
    );
  };

  const SourceCode = () => {
    return (
      <SyntaxHighlighter language={language}>
        {sourceCode}
      </SyntaxHighlighter>
    );
  };

  const getAllRuns = () => {
    fetch("/api/admin/all_runs")
      .then((response) => response.json())
      .then((data) => setRuns(data.reverse()));
  };

  useEffect(() => {
    getAllRuns();
  }, []);

  return (
    <div> 
      <h1>All Runs</h1>
      <RunTable runs={runs} />
      {sourceCode && <SourceCode />}
    </div>
  );
};

export default AllRuns;
