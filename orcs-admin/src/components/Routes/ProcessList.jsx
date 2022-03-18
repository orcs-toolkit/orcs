import moment from "moment";
import axios from "axios";
import { useLocation } from "react-router-dom";
import React, { useState, useEffect, Fragment } from "react";
import StyledTable from "../utilities/Table";
import "rsuite/dist/rsuite.min.css";
import { Table } from "rsuite";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Button } from "react-bootstrap";

const ProcessList = ({ socket }) => {
  const [dataList, setDataList] = useState([]);
  const [role, setrole] = useState("");
  const [loading, setLoading] = useState(false);
  const url = useLocation().search;
  const paramValue = new URLSearchParams(url).get("macA") || "";

  useEffect(() => {
    const getDataListener = (data) => {
      if (dataList.length === 0) {
        let tempMacA = Object.keys(data)[0];
        if (tempMacA === paramValue) {
          setrole(data[tempMacA].role || "");
          setDataList(
            data[tempMacA]["processData"]["list"].filter(
              (p) => p.state !== "unknown"
            ) || []
          );
        }
      }
    };
    socket.on("data", getDataListener);
  }, [socket]);

  const addSingleProcess = (role, processName) => {
    if (role !== "") {
      setLoading(true);
      axios
        .post(`http://localhost:4001/api/updateSinglePolicy`, {
          role,
          processName,
        })
        .then((res) => {
          toast("Successfully added process to policy list");
          setLoading(false);
          console.log("Success", res);
        })
        .catch((err) => {
          toast("Something went wrong");
          setLoading(false);
          console.log(err);
        });
    }
  };

  return (
    <Fragment>
      <StyledTable
        title={"Process List"}
        dataList={dataList}
        loading={dataList.length === 0}
        tableBody={
          <>
            <Table.Column width={120} align="center" fixed>
              <Table.HeaderCell>PID</Table.HeaderCell>
              <Table.Cell dataKey="pid" />
            </Table.Column>
            <Table.Column width={170} fixed>
              <Table.HeaderCell>Process Name</Table.HeaderCell>
              <Table.Cell dataKey="name" />
            </Table.Column>
            <Table.Column width={140}>
              <Table.HeaderCell>Duration</Table.HeaderCell>
              <Table.Cell>
                {(rowData) =>
                  `${moment().diff(rowData.started, "hours")} hours`
                }
              </Table.Cell>
            </Table.Column>
            <Table.Column width={150}>
              <Table.HeaderCell>Status</Table.HeaderCell>
              <Table.Cell>
                {(rowData) => (
                  <div className="d-flex align-items-center">
                    <i
                      className={`bx bxs-circle text-${
                        rowData.state === "sleeping"
                          ? "warning"
                          : rowData.state === "running"
                          ? "success"
                          : "dark"
                      } mr-1`}
                    ></i>{" "}
                    {rowData.state[0].toUpperCase() +
                      rowData.state.substring(1)}
                  </div>
                )}
              </Table.Cell>
            </Table.Column>

            <Table.Column width={120}>
              <Table.HeaderCell>Add to Ban List</Table.HeaderCell>
              <Table.Cell>
                {(rowData) => (
                  <Button
                    disabled={loading}
                    onClick={() => addSingleProcess(role, rowData.name)}
                    className="btn btn-danger ml-3 shadow btn-xs sharp"
                  >
                    <i className="bx bx-plus fs-15"></i>
                  </Button>
                )}
              </Table.Cell>
            </Table.Column>
          </>
        }
      />
    </Fragment>
  );
};

export default ProcessList;