import React from "react";
import { Menu, X } from "lucide-react";
import Button from "@/components/Button";

const MobileMenuToggleButton = ({ isOpen, onClick }) => {
  return (
    <Button
      onClick={onClick}
      className="navigation__mobile-button"
      aria-label={isOpen ? "Закрыть меню" : "Открыть меню"}
    >
      {isOpen ? (
        <X className="navigation__mobile-icon" />
      ) : (
        <Menu className="navigation__mobile-icon" />
      )}
    </Button>
  );
};

export default MobileMenuToggleButton;
