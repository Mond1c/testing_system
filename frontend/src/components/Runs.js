import React, { useEffect, useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import styled from "styled-components";

const TableWrapper = styled.table`
  border: 1px solid black;
  width: 100%;
`;

const TableHeaderWrapper = styled.th`
  border: 1px solid black;
`;

const TableElementWrapper = styled.td`
  border: 1px solid black;
  text-align: center;
`;


const Runs = () => {
    const [runs, setRuns] = useState([]);
    const [sourceCode, setSourceCode] = useState(undefined);
    const [language, setLanguage] = useState("");

    const getRuns = () => {
        fetch("/api/runs")
            .then(response => response.json())
            .then(response => setRuns(response));
    };

    useEffect(() => {
        getRuns();
    }, []);

    const SourceCode = () => {
        return (
          <SyntaxHighlighter language={language}>
            {sourceCode}
          </SyntaxHighlighter>
        );
      };

    return (
        <div>
            <h1>Runs</h1>
            <TableWrapper>
                <tr>
                    <TableHeaderWrapper>Problem</TableHeaderWrapper>
                    <TableHeaderWrapper>Result</TableHeaderWrapper>
                    <TableHeaderWrapper>Time</TableHeaderWrapper>
                    <TableHeaderWrapper>Source Code</TableHeaderWrapper>
                </tr>
            {
                runs.reverse().map(run => {
                    const viewSourceCode = () => {
                        fetch("/api/source_code?run_id="+run.run_id)
                            .then(response => response.text())
                            .then(data => {
                                setSourceCode(data);
                                setLanguage(run.language);
                            });
                    }
                    return (<tr>
                            <TableElementWrapper>{run.problem}</TableElementWrapper>
                            <TableElementWrapper>{run.result}</TableElementWrapper>
                            <TableElementWrapper>{run.time}</TableElementWrapper>
                            <TableElementWrapper>
                                <button onClick={viewSourceCode}>View</button>
                            </TableElementWrapper>
                        </tr>
                    )
                })
            }
            </TableWrapper>
            {sourceCode && <SourceCode/>}
        </div>
    );
};

export default Runs;
