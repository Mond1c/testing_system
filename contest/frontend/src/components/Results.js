import React, { useEffect, useState } from "react";
import styled from "styled-components";

const Wrapper = styled.div`
  margin: 20px;
`;

const Title = styled.h1`
  font-size: 24px;
  margin-bottom: 10px;
`;

const TableWrapper = styled.table`
  border-collapse: collapse;
  width: 100%;
`;

const TableHeaderWrapper = styled.th`
  border: 1px solid black;
  padding: 8px;
  background-color: #f2f2f2;
  font-weight: bold;
`;

const TableElementWrapper = styled.td`
  border: 1px solid black;
  padding: 8px;
  text-align: center;
`;

const Results = () => {
  const [results, setResults] = useState([]);
  const [problems, setProblems] = useState([]);

  const getResults = () => {
    console.log("Results getting");
    fetch("/api/results")
      .then((response) => response.json())
      .then((response) => {
        setProblems(response.problems);
        const res = Object.values(response.contestants);
        console.log(res);
        res.sort((a, b) => {
          if (a.points === b.points) {
            return a.penalty - b.penalty;
          }
          return b.points - a.points;
        });
        setResults(res);
      });
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
    console.log(contestants);
    return contestants.map((contestant) => {
      return (
        <tr key={contestant.id}>
          <TableElementWrapper>{contestant.id}</TableElementWrapper>
          <TableElementWrapper>{contestant.name}</TableElementWrapper>
          {problems.map((problem) => {
            return getBestResult(contestant, problem);
          })}
          <TableElementWrapper>{contestant.points}</TableElementWrapper>
          <TableElementWrapper>{contestant.penalty}</TableElementWrapper>
        </tr>
      );
    });
  };

  return (
    <Wrapper>
      <Title>Results</Title>
      <TableWrapper>
        <thead>
          <tr>
            <TableHeaderWrapper>Id</TableHeaderWrapper>
            <TableHeaderWrapper>Name</TableHeaderWrapper>
            {problems.map((problem) => {
              console.log(results);
              return (
                <TableHeaderWrapper key={problem}>{problem}</TableHeaderWrapper>
              );
            })}
            <TableHeaderWrapper>Score</TableHeaderWrapper>
            <TableHeaderWrapper>Penalty</TableHeaderWrapper>
          </tr>
        </thead>
        <tbody>
          {results !== undefined ? getContestants(results) : null}
        </tbody>
      </TableWrapper>
    </Wrapper>
  );
};

export default Results;
