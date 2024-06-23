import React, { useEffect, useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import styled from "styled-components";

const Wrapper = styled.div`
    padding: 20px;
`;

const Title = styled.h1`
    font-size: 24px;
    margin-bottom: 20px;
`;

const Table = styled.table`
    border-collapse: collapse;
    width: 100%;
`;

const TableHeader = styled.th`
    border: 1px solid black;
    padding: 10px;
    background-color: #f2f2f2;
    font-weight: bold;
`;

const TableCell = styled.td`
    border: 1px solid black;
    padding: 10px;
    text-align: center;
`;

const Button = styled.button`
    padding: 5px 10px;
    background-color: #007bff;
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
`;

const Runs = () => {
    const [runs, setRuns] = useState([]);
    const [sourceCode, setSourceCode] = useState(undefined);
    const [language, setLanguage] = useState("");

    const getRuns = () => {
        fetch("/api/runs")
          .then((response) => response.json())
          .then((data) => setRuns(data.reverse()));
    };

    useEffect(() => {
        getRuns();
    }, []);

    const viewSourceCode = (run) => {
        fetch("/api/source_code?run_id=" + run.run_id)
            .then((response) => response.text())
            .then((data) => {
                setSourceCode(data);
                setLanguage(run.language);
            });
    };

    const SourceCode = () => {
        return (
            <SyntaxHighlighter language={language}>
                {sourceCode}
            </SyntaxHighlighter>
        );
    };

    return (
        <Wrapper>
            <Title>Runs</Title>
            <Table>
                <thead>
                    <tr>
                        <TableHeader>Problem</TableHeader>
                        <TableHeader>Result</TableHeader>
                        <TableHeader>Time</TableHeader>
                        <TableHeader>Source Code</TableHeader>
                    </tr>
                </thead>
                <tbody>
                    {runs.map((run) => (
                        <tr key={run.run_id}>
                            <TableCell>{run.problem}</TableCell>
                            <TableCell>{run.result}</TableCell>
                            <TableCell>{run.time}</TableCell>
                            <TableCell>
                                <Button onClick={() => viewSourceCode(run)}>View</Button>
                            </TableCell>
                        </tr>
                    ))}
                </tbody>
            </Table>
            {sourceCode && <SourceCode />}
        </Wrapper>
    );
};

export default Runs;
