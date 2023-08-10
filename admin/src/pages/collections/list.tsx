import { BaseRecord, IResourceComponentsProps, useMany } from "@refinedev/core";
import { AntdListInferencer } from "@refinedev/inferencer/antd";

import React from "react";
import {
  useTable,
  List,
  EditButton,
  ShowButton,
  TextField,
  ImageField,
} from "@refinedev/antd";
import { Table, Space } from "antd";

export const CollectionsList: React.FC<IResourceComponentsProps> = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  return (
    <List>
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
        Collections is the section where you can create a collection for your
        SBT.
      </Space>
      <Table size="large" {...tableProps} rowKey="id">
        <Table.Column
          dataIndex="raw_address"
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
          title="Raw Address"
        />
        <Table.Column
          dataIndex="friendly_address"
          title="Friendly Address"
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
        />
        <Table.Column dataIndex="next_item_index" title="Next Item Index" />
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
          dataIndex="content_uri"
          title="Content Uri"
        />
        <Table.Column dataIndex="name" title="Name" />
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
          dataIndex="description"
          title="Description"
        />
        <Table.Column
          dataIndex="image"
          render={(value) => {
            return <ImageField value={value} style={{ maxWidth: 100 }} />;
          }}
          title="Image"
        ></Table.Column>
        <Table.Column dataIndex="default_weight" title="Default Rating Points" />
        <Table.Column dataIndex="version" title="Version" />
        <Table.Column
          title="Actions"
          dataIndex="actions"
          render={(_, record: BaseRecord) => (
            <Space>
              <EditButton hideText size="small" recordItemId={record.id} />
              <ShowButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};
