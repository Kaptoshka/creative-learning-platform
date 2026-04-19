import React from "react";

const TaskInstructions = ({ prompt, example, instructions }) => {
  return (
    <>
      <p className="instructions">
        {prompt ||
          "Составьте предложение, где все слова будут начинаться на одну букву."}
      </p>

      {example && (
        <div className="example">
          <strong>Пример:</strong>
          <p>{example}</p>
        </div>
      )}

      {instructions && instructions.length > 0 && (
        <div className="instructions">
          <strong>Инструкции:</strong>
          <ol>
            {instructions.map((instruction, index) => (
              <li key={index}>{instruction}</li>
            ))}
          </ol>
        </div>
      )}
    </>
  );
};

export default TaskInstructions;
