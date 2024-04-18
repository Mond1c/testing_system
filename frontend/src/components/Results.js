import React, { useEffect, useState } from "react";
import styled from "styled-components";

const TableWrapper = styled.table`
  border: 1px solid black;
`;

const TableHeaderWrapper = styled.th`
  border: 1px solid black;
`;

const TableElementWrapper = styled.td`
  border: 1px solid black;
  text-align: center;
`;

const Results = () => {
  const [results, setResults] = useState({});

  const getResults = () => {
    console.log("Results getting");
    fetch("/api/results")
      .then((response) => response.json())
      .then((response) => setResults(response));
  };

  useEffect(() => getResults(), []);

  const getBestResult = (contestant, problem) => {
    if (
      contestant.results != null &&
      contestant.results[problem] !== undefined
    ) {
      const result = contestant.results[problem].result.result;
      if (result === 1) {
        const penalty = contestant.additionalPenalty[problem] / 20;
        return (
          <TableElementWrapper>
            +{penalty > 0 ? penalty : ""}
          </TableElementWrapper>
        );
      }
      return <TableElementWrapper>-</TableElementWrapper>;
    }
    return <TableElementWrapper>.</TableElementWrapper>;
  };

  const getContestants = (contestants) => {
    return Object.entries(contestants)?.map(([_, contestant]) => {
      return (
        <tr>
          <TableElementWrapper>{contestant.id}</TableElementWrapper>
          <TableElementWrapper>{contestant.name}</TableElementWrapper>
          {results.problems.map((problem) => {
            return getBestResult(contestant, problem);
          })}
          <TableElementWrapper>{contestant.points}</TableElementWrapper>
          <TableElementWrapper>{contestant.penalty}</TableElementWrapper>
        </tr>
      );
    });
  };

  return (
    <div>
      <h1>Results</h1>
      <TableWrapper style={{ width: "100%" }}>
        <tr>
          <TableHeaderWrapper>Id</TableHeaderWrapper>
          <TableHeaderWrapper>Name</TableHeaderWrapper>
          {results?.problems?.map((problem) => {
            console.log(results);
            return <TableHeaderWrapper>{problem}</TableHeaderWrapper>;
          })}
          <TableHeaderWrapper>Score</TableHeaderWrapper>
          <TableHeaderWrapper>Penalty</TableHeaderWrapper>
        </tr>
        {results?.contestants !== undefined
          ? getContestants(results?.contestants)
          : ""}
      </TableWrapper>
    </div>
  );
};

export default Results;
