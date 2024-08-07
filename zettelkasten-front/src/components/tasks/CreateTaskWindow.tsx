import React, { useState, KeyboardEvent } from "react";
import { saveNewTask } from "../..//api/tasks";

import { Task, emptyTask } from "../../models/Task";
import { Card, PartialCard } from "../..//models/Card";
import { BacklinkInput } from "../cards/BacklinkInput";
import { TaskDateDisplay } from "./TaskDateDisplay";
import { Input } from "../Input";

interface CreateTaskWindowProps {
  currentCard: Card | PartialCard | null;
  setRefresh: (refresh: boolean) => void;
  setShowTaskWindow: (showTaskWindow: boolean) => void;
}

export function CreateTaskWindow({
  currentCard,
  setRefresh,
  setShowTaskWindow,
}: CreateTaskWindowProps) {
  const [newTask, setNewTask] = useState<Task>(emptyTask);
  const [selectedCard, setSelectedCard] = useState<PartialCard | null>(null);
  async function handleSaveTask() {
    let response;

    console.log("create card", currentCard);

    let task = newTask;
    if (currentCard) {
      task = { ...task, card_pk: currentCard.id };
    }
    response = await saveNewTask(task);
    if (!("error" in response)) {
      setRefresh(true);
      setShowTaskWindow(false);
      setSelectedCard(null);
      let date = newTask.scheduled_date;
      setNewTask({ ...emptyTask, scheduled_date: date });
      if (currentCard) {
        setNewTask({ ...emptyTask, card_pk: currentCard.id });
      }
    }
  }
  function handleBacklink(card: PartialCard) {
    setSelectedCard(card);
    setNewTask({ ...newTask, card_pk: card.id });
  }

  return (
    <div
      className="create-task-popup-overlay"
      onClick={() => setShowTaskWindow(false)}
    >
      <div
        className="create-task-popup-content"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="create-task-window">
          <div className="create-task-window-top">
            <Input
              label="New Task"
              value={newTask.title}
              onChange={(e) =>
                setNewTask({ ...newTask, title: e.target.value })
              }
              onKeyPress={(event: React.KeyboardEvent<HTMLInputElement>) => {
                if (event.key === "Enter") {
                  handleSaveTask();
                }
              }}
              placeholder="Enter new task"
              className="w-full"
              autoFocus
            />
          </div>
          <div className="create-task-window-bottom">
            <div className="create-task-window-bottom-left">
              {!currentCard && <BacklinkInput addBacklink={handleBacklink} />}
            </div>
            <div className="crate-task-window-bottom-middle">
              {selectedCard && (
                <span style={{ fontWeight: "bold", color: "blue" }}>
                  [{selectedCard?.card_id}]
                </span>
              )}
              <TaskDateDisplay
                task={newTask}
                setTask={setNewTask}
                setRefresh={setRefresh}
                saveOnChange={false}
              />
            </div>
            <div className="create-task-window-bottom-right">
              <button onClick={handleSaveTask}>Save</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
