import React, { useState, useEffect } from "react";
import styled from "styled-components";

const DurationContainer = styled.div`
    background-color: #f2f2f2;
    padding: 10px;
    border-radius: 5px;
    margin-bottom: 10px;
`;

const DurationText = styled.p`
    font-size: 16px;
    color: #333;
    margin-bottom: 5px;
`;

const Duration = ({ startTime, duration }) => {
    return (
        <DurationContainer>
            <DurationText>Started at: {new Date(startTime).toLocaleString()}</DurationText>
            <DurationText>Duration: {Math.floor(Date.now() - new Date(startTime)) / 1000 } seconds of {duration} seconds</DurationText>
        </DurationContainer>
    );
}

const MainContainer = styled.div`
    background-color: #fff;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px;
`;

const MainHeading = styled.h1`
    font-size: 24px;
    color: #333;
    margin-bottom: 10px;
`;

const Main = () => {
    const [startTime, setStartTime] = useState(Date.now());
    const [duration, setDuration] = useState(0);
    const [username, setUsername] = useState("");

    const getStartTime = () => {
        fetch("/api/startTime")
        .then(response => response.json())
        .then(response => {
            setStartTime(response.startTime);
            setDuration(response.duration);
        });
    };

    const getUsername = () => {
        fetch("/api/me")
        .then(response => response.json())
        .then(response => setUsername(response.username));
    }

    useEffect(() => {
        getStartTime();
        getUsername();
    }, []);

    return (
        <MainContainer>
            <MainHeading>Welcome, {username}</MainHeading>
            <Duration startTime={startTime} duration={duration} />
        </MainContainer>
    );
};

export default Main;