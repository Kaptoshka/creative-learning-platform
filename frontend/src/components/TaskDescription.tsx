import React from "react";

const TaskDescription = ({ instructions, example }) => {
  return (
    <>
      <p className="example">{instructions}</p>
      {example && (
        <p className="example">
          <strong>Например:</strong>
          <br />
          {example}
        </p>
      )}
    </>
  );
};

export default TaskDescription;
