import React, { useState, KeyboardEvent } from "react";
import { deleteTask, saveExistingTask } from "../../api/tasks";
import { getTomorrow } from "../../utils/dates";

import { TaskDateDisplay } from "./TaskDateDisplay";
import { Task } from "../../models/Task";
import { Link } from "react-router-dom";
import { PartialCard } from "../../models/Card";
import { BacklinkInput } from "../cards/BacklinkInput";
import { linkifyWithDefaultOptions } from "../../utils/strings";
import { TaskClosedIcon } from "../../assets/icons/TaskClosedIcon";
import { TaskOpenIcon } from "../../assets/icons/TaskOpenIcon";
import { TaskTagDisplay } from "../../components/tasks/TaskTagDisplay";
import { AddTagMenu } from "../../components/tasks/AddTagMenu";

interface TaskListItemProps {
  task: Task;
  setRefresh: (refresh: boolean) => void;
  onTagClick: (tag: string) => void;
}

export function TaskListItem({
  task,
  setRefresh,
  onTagClick,
}: TaskListItemProps) {
  const [editTitle, setEditTitle] = useState<boolean>(false);
  const [newTitle, setNewTitle] = useState<string>("");
  const [showCardLink, setShowCardLink] = useState<boolean>(false);
  const [showMenu, setShowMenu] = useState<boolean>(false);
  const [showTagMenu, setShowTagMenu] = useState<boolean>(false);

  async function handleTitleClick() {
    setNewTitle(task.title);
    setEditTitle(true);
  }

  function toggleCardLink() {
    setShowCardLink(!showCardLink);
    setShowMenu(false);
  }
  async function handleCardUnlink() {
    let editedTask = { ...task, card_pk: 0 };
    let response = await saveExistingTask(editedTask);
    if (!("error" in response)) {
      setRefresh(true);
    }
    setShowMenu(false);
  }

  async function handleBacklink(card: PartialCard) {
    let editedTask = { ...task, card_pk: card.id };
    let response = await saveExistingTask(editedTask);
    if (!("error" in response)) {
      setRefresh(true);
      setShowCardLink(false);
    }
  }

  async function handleTitleEdit() {
    let editedTask = { ...task, title: newTitle };
    let response = await saveExistingTask(editedTask);
    if (!("error" in response)) {
      setRefresh(true);
      setEditTitle(false);
      setNewTitle("");
    }
  }

  async function handleToggleComplete() {
    let editedTask = { ...task, is_complete: task.is_complete ? false : true };
    let response = await saveExistingTask(editedTask);
    if (!("error" in response)) {
      setRefresh(true);
    }
  }
  async function handleAddTagClick() {
    setShowMenu(false);
    setShowTagMenu(true);
  }
  function toggleMenu() {
    if (showTagMenu) {
      setShowTagMenu(false);
    }
    setShowMenu(!showMenu);
  }
  async function handleDelete() {
    let _ = await deleteTask(task.id);
    setRefresh(true);
    setShowMenu(false);
  }
  return (
    <div className="task-list-item">
      <div className="task-list-item-checkbox">
        <span onClick={handleToggleComplete}>
          {task.is_complete ? <TaskClosedIcon /> : <TaskOpenIcon />}
        </span>
      </div>
      <div className="task-list-item-middle-container">
        <div className="task-list-item-title">
          {editTitle ? (
            <input
              className="task-list-item-title-input"
              value={newTitle}
              onChange={(e) => setNewTitle(e.target.value)}
              onKeyPress={(event: KeyboardEvent<HTMLInputElement>) => {
                if (event.key === "Enter") {
                  handleTitleEdit();
                }
              }}
            />
          ) : (
            <span
              onClick={handleTitleClick}
              className={task.is_complete ? "task-completed" : "task-title"}
              dangerouslySetInnerHTML={{
                __html: linkifyWithDefaultOptions(task.title),
              }}
            />
          )}
        </div>
        <div className="task-list-item-details inline-block">
          <TaskDateDisplay
            task={task}
            setTask={(task: Task) => {}}
            setRefresh={setRefresh}
            saveOnChange={true}
          />
          <TaskTagDisplay task={task} onTagClick={onTagClick} />
        </div>
      </div>
      <div className="task-list-item-card">
        {task.card && task.card.id > 0 && (
          <Link
            to={`/app/card/${task.card.id}`}
            style={{ textDecoration: "none", color: "inherit" }}
          >
            <span className="card-id">[{task.card.card_id}]</span>
          </Link>
        )}
        {!task.card ||
          (task.card.id == 0 && (
            <div>
              {showCardLink && <BacklinkInput addBacklink={handleBacklink} />}
            </div>
          ))}
      </div>
      <div className="dropdown">
        <button onClick={toggleMenu} className="menu-button">
          ⋮
        </button>
        {showMenu && (
          <div className="popup-menu">
            {task.card_pk === 0 ? (
              <button onClick={() => toggleCardLink()}>Link Card</button>
            ) : (
              <button onClick={() => handleCardUnlink()}>Unlink Card</button>
            )}
            <button onClick={() => handleAddTagClick()}>Add Tag</button>
            <button onClick={() => handleDelete()}>Delete</button>
          </div>
        )}
        {showTagMenu && (
          <div className="popup-menu">
            <AddTagMenu
              task={task}
              setRefresh={setRefresh}
              setShowTagMenu={setShowTagMenu}
            />
          </div>
        )}
      </div>
    </div>
  );
}
