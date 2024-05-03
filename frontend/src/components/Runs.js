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

const Runs = () => {
    const [runs, setRuns] = useState([]);
    const [username, setUsername] = useState(undefined);

    const getResult = (r) => {
        const result = Results[r.result];
        if (r.number > -1) {
          return result + " on test " + (r.number + 1);
        } else {
          return result;
        }
    }

    const getUsername = () => {
        return fetch("/api/me")
            .then((response) => response.json());
    };

    const getRuns = (name) => {
        fetch("/api/runs?name=" + name)
            .then(response => response.json())
            .then(response => setRuns(response));
    };

    useEffect(() => {
        getUsername().then(response => {
            setUsername(response.username);
            getRuns(response.username);
        });
    }, []);

    return (
        <div>
            <h1>Hello, {username !== undefined ? username : ""}</h1>
            {
                runs.map(run => {
                    return (<div>
                        <p>problem: {run.problem}, result: {getResult(run.result)}, time: {run.time}</p>
                        <br />
                        </div>
                    )
                })
            }
        </div>
    );
};

export default Runs;
