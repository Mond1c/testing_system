import React, { useEffect, useState } from "react";

const Results = () => {
  const [results, setResults] = useState({});

  const getResults = () => {
    console.log("Results getting");
    fetch("/api/results")
      .then((response) => response.json())
      .then((response) => setResults(response));
  };

  useEffect(() => getResults(), []);

  return (
    <div>
      <h1>Results</h1>
      {JSON.stringify(results)}
    </div>
  );
};

export default Results;
