import React, { useEffect, useState } from "react";
import styled from "styled-components";

const Results = {
    0: "Undefined behaviour",
    1: "OK",
    2: "Compile error",
    3: "Runtime error",
    4: "Time limit",
    5: "Memory limit",
    6: "Wrong answer"
};

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

    const getResult = (r) => {
        const result = Results[r.result];
        if (r.number > -1) {
          return result + " on test " + (r.number + 1);
        } else {
          return result;
        }
    }

    const getRuns = () => {
        fetch("/api/runs")
            .then(response => response.json())
            .then(response => setRuns(response));
    };

    useEffect(() => {
        getRuns();
    }, []);

    return (
        <div>
            <h1>Runs</h1>
            <TableWrapper>
                <tr>
                    <TableHeaderWrapper>Problem</TableHeaderWrapper>
                    <TableHeaderWrapper>Result</TableHeaderWrapper>
                    <TableHeaderWrapper>Time</TableHeaderWrapper>
                </tr>
            {
                runs.reverse().map(run => {
                    return (<tr>
                            <TableElementWrapper>{run.problem}</TableElementWrapper>
                            <TableElementWrapper>{getResult(run.result)}</TableElementWrapper>
                            <TableElementWrapper>{run.time}</TableElementWrapper>
                        </tr>
                    )
                })
            }
            </TableWrapper>
        </div>
    );
};

export default Runs;
