import React, { useState, useEffect } from "react";

const Duration = ({ startTime, duration }) => {
    return (
        <div>
            <p>Started at: {new Date(startTime).toLocaleString()}</p>
            <p>Duration: {Math.floor(Date.now() - new Date(startTime)) / 1000 } seconds of {duration} seconds</p>
        </div>
    );
}

const Main = () => {
    const [startTime, setStartTime] = useState(Date.now());
    const [duration, setDuration] = useState(0);

    const getStartTime = () => {
        fetch("/api/startTime")
        .then(response => response.json())
        .then(response => {
            setStartTime(response.startTime);
            setDuration(response.duration);
        });
    }

    useEffect(() => {
        getStartTime();
    }, []);

    return (
        <div>
            <h1>Home</h1>
            <Duration startTime={startTime} duration={duration} />
        </div>
    );
};

export default Main;
