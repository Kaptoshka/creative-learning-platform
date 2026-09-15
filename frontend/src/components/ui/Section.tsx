import React from "react";

interface SectionProps {
  title: string;
  icon?: React.ComponentType<{ size?: number }>;
  action?: {
    label: string;
    onClick: () => void;
  };
  children: React.ReactNode;
}

const Section = ({ title, icon: Icon, action, children }: SectionProps) => {
  return (
    <div className="section">
      <div className="section-header">
        {Icon && <Icon size={20} />}
        <h2>{title}</h2>
        {action && (
          <button onClick={action.onClick} className="section-action">
            {action.label}
          </button>
        )}
      </div>
      {children}
    </div>
  );
};

export default Section;
