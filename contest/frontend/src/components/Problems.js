import React, { useEffect, useState } from 'react';
import styled from 'styled-components';

const Container = styled.div`
    padding: 20px;
    background-color: #f9f9f9;
`;

const Table = styled.table`
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 20px;
`;

const Th = styled.th`
    padding: 12px 8px;
    text-align: left;
    border-bottom: 2px solid #ddd;
    background-color: #f2f2f2;
`;

const Td = styled.td`
    padding: 8px;
    border-bottom: 1px solid #ddd;
    vertical-align: top; /* Ensure vertical alignment for nested table */
`;

const Problems = () => {
    const [runs, setRuns] = useState([]);
    const [problems, setProblems] = useState([]);

    const getRuns = () => {
        fetch('/api/runs')
            .then((response) => response.json())
            .then((data) => setRuns(data.reverse()));
    };

    const getProblems = () => {
        fetch('/api/problems')
            .then((response) => response.json())
            .then((data) => setProblems(data));
    };

    useEffect(() => {
        getRuns();
        getProblems();
    }, []);

    return (
        <Container>
            <Table>
                <thead>
                    <tr>
                        <Th>Problem</Th>
                        <Th>Runs</Th>
                    </tr>
                </thead>
                <tbody>
                    {problems.map((problem) => (
                        <tr key={problem}>
                            <Td>{problem}</Td>
                            <Td>
                                <Table>
                                    <thead>
                                        <tr>
                                            <Th>Id</Th>
                                            <Th>Result</Th>
                                            <Th>Time</Th>
                                            <Th>Language</Th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {runs
                                            .filter((run) => run.problem === problem)
                                            .map((run) => (
                                                <tr key={run.run_id}>
                                                    <Td>{run.run_id}</Td>
                                                    <Td>{run.result}</Td>
                                                    <Td>{run.time}</Td>
                                                    <Td>{run.language}</Td>
                                                </tr>
                                            ))}
                                    </tbody>
                                </Table>
                            </Td>
                        </tr>
                    ))}
                </tbody>
            </Table>
        </Container>
    );
};

export default Problems;
