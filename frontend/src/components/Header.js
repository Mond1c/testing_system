import React from 'react';
import styled from 'styled-components';

const HeaderContainer = styled.header`
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 100px;
    background-color: #f2f2f2;
    padding: 0 20px;
`;

const Nav = styled.nav`
    display: flex;
    align-items: center;
`;

const NavLink = styled.a`
    margin: 0 10px;
    color: #333;
    text-decoration: none;
    font-size: 18px;
    font-weight: bold;

    &:hover {
        color: #666;
    }
`;

const Header = () => {
    return (
        <HeaderContainer>
            <Nav>
                <NavLink href="/">Home</NavLink>
                <NavLink href="/results">Results</NavLink>
                <NavLink href="/runs">Runs</NavLink>
                <NavLink href="/problems">Problems</NavLink>
                <NavLink href="/upload">Upload Solution</NavLink>
                <NavLink href="/condition">Condition</NavLink>
            </Nav>
        </HeaderContainer>
    );
};

export default Header;