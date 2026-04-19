import React from "react";
import { NavLink } from "react-router-dom";
import Button from "@/components/Button";

const NavigationLink = ({ to, end = false, onClick, icon: Icon, children }) => {
  const buttonVariant = onClick ? "mobile-navigation" : "navigation";

  return (
    <NavLink to={to} end={end} viewTransition onClick={onClick}>
      {({ isActive }) => (
        <Button
          variant={buttonVariant}
          isActive={isActive && buttonVariant === "navigation"}
          icon={Icon}
        >
          {children}
        </Button>
      )}
    </NavLink>
  );
};

export default NavigationLink;
