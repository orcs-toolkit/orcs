import React, { useState, useEffect } from "react";
import api from "../../../utils/api";
import { Button, Modal, Form } from "react-bootstrap";
import AsyncSelect from "react-select/async";

const AddProcessModal = ({ visible, toggle, processName, socket }) => {
  const [roles, setRoles] = useState([]);
  const [addToFavorite, setAddToFavorite] = useState([]);
  const [selectedRoles, setSelectedRoles] = useState([]);
  const [error, setError] = useState(false);

  const handleSubmit = () => {
    console.log("srjjjj", selectedRoles);
    if (selectedRoles !== [] || selectedRoles[0].value === "All") {
      setError(false);
      api
        .post(`http://localhost:4001/policy/addSingleProcess`, {
          processName,
          ...(addToFavorite && { addToFavorite }),
          role: selectedRoles.map((b) => b.value),
        })
        .then((res) => {
          console.log(socket);
          socket.emit("updated:Ban", {
            selectedRoles,
            status: "updated",
          });
          console.log("Success", res);
          setSelectedRoles([]);
          setError(false);
          toggle();
        })
        .catch((err) => {
          setError(true);
          console.log(err);
        });
    } else {
      setError(true);
    }
  };

  const fetchRoles = () => {
    api
      .get("http://localhost:4001/policy/getRoleWisePolicy")
      .then((res) => {
        setRoles(res.data.map((m) => ({ value: m._id, label: m._id })));
      })
      .catch((err) => {
        console.log(err);
      });
  };

  useEffect(() => {
    fetchRoles();
  }, []);

  return (
    <Modal className="fade" show={visible}>
      <Modal.Header>
        <Modal.Title>Select Roles</Modal.Title>
        <Button variant="" className="close" onClick={toggle}>
          <span>&times;</span>
        </Button>
      </Modal.Header>
      <Modal.Body>
        <Form.Group className="my-3">
          <Form.Label>Selected Process</Form.Label>
          <input
            type="text"
            className="w-100 p-2"
            defaultValue={processName}
            name="Role"
            disabled={true}
            placeholder="Selected Process"
          />
        </Form.Group>
        <Form.Group className="my-3">
          <Form.Label>Select Roles</Form.Label>
          <AsyncSelect
            styles={{
              menuPortal: (base) => ({
                ...base,
                zIndex: 9999,
              }),
            }}
            isMulti={selectedRoles[0]?.value !== "All"}
            defaultOptions={[{ value: "All", label: "All" }, ...roles]}
            value={selectedRoles}
            onChange={(newValue) => {
              setSelectedRoles(newValue);
            }}
            placeholder={"Select Roles"}
          />
        </Form.Group>
        <Form.Group className="my-3">
          <Form.Label>Add To Favorite List</Form.Label>
          <AsyncSelect
            styles={{
              menuPortal: (base) => ({
                ...base,
                zIndex: 9999,
              }),
            }}
            defaultValue={{ value: false, label: "No" }}
            defaultOptions={[
              { value: true, label: "Yes" },
              { value: false, label: "No" },
            ]}
            onChange={(newValue) => {
              setAddToFavorite(newValue.value);
            }}
            placeholder={"Add To Favorite List"}
          />
        </Form.Group>
        {error && <h3 className="text-danger">Some error occured</h3>}
      </Modal.Body>
      <Modal.Footer>
        <Button onClick={toggle} variant="danger rounded">
          Close
        </Button>
        <Button variant="primary  rounded" onClick={handleSubmit}>
          Save changes
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default AddProcessModal;