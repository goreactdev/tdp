import React from "react";
import {
  IResourceComponentsProps,
  BaseRecord,
  HttpError,
} from "@refinedev/core";
import {
  useTable,
  List,
  EditButton,
  ShowButton,
  SaveButton,
  TextField,
  ImageField,
} from "@refinedev/antd";
import { Table, Space, Form, Input } from "antd";
import { CreateButton } from "../../components/buttons/create";

export const SBTTokensList: React.FC<IResourceComponentsProps> = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  return (
    <List  headerButtons={<CreateButton />}>
      <Space
        style={{
          marginBottom: "20px",
          maxWidth: "25rem",
          backgroundColor: "#fff",
          borderRadius: "10px",
          padding: "10px",
          fontWeight: "600",
          fontSize: "14px",
        }}
        direction="vertical"
        size="small"
      >
        Minted NFTs is a section where you can mint NFTs to users. All NFTs are SBT tokens. Default rating is provided from the collection. You can rewrite rating for selected SBTs.
        <span
          style={{
            color: "red",
          }}
        >
          You can mint NFTs only to the collection YOU own.<br />
          1. Create a collection<br />
          2. Mint NFTs to users
        </span>

      </Space>

      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="Id" />
        <Table.Column
          render={(value) => {
            return (
              <TextField
                value={
                  value.length < 10
                    ? value
                    : value.slice(0, 10) + "..." + value.slice(-4)
                }
              />
            );
          }}
          dataIndex="raw_address"
          title="Raw addr"
        />
        <Table.Column
          render={(value) => {
            return (
              <TextField
                value={
                  value.length < 10
                    ? value
                    : value.slice(0, 10) + "..." + value.slice(-4)
                }
              />
            );
          }}
          dataIndex="friendly_address"
          title="Friendly addr"
        />

        <Table.Column dataIndex="name" title="Name" />
        <Table.Column
          dataIndex="description"
          render={(value) => {
            return (
              <TextField
                value={
                  value.length < 10
                    ? value
                    : value.slice(0, 10) + "..." + value.slice(-4)
                }
              />
            );
          }}
          title="Description"
        />
        <Table.Column
          dataIndex="image"
          render={(value) => {
            return <ImageField value={value} style={{ maxWidth: 100 }} />;
          }}
          title="Image"
        ></Table.Column>
        <Table.Column dataIndex="weight" title="Rating Points" />
        <Table.Column dataIndex="index" title="Index" />
        <Table.Column dataIndex="version" title="Version" />
        <Table.Column
          title="Actions"
          dataIndex="actions"
          render={(_, record: BaseRecord) => (
            <Space>
              <ShowButton hideText size="small" recordItemId={record.id} />
              <EditButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};
