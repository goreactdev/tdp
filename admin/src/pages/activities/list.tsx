import React from "react";
import { IResourceComponentsProps, BaseRecord } from "@refinedev/core";
import {
  useTable,
  List,
  EditButton,
  ShowButton,
  MarkdownField,
} from "@refinedev/antd";
import { Table, Space } from "antd";
import { CheckCircleFilled, CheckOutlined } from "@ant-design/icons";
import { Link } from "react-router-dom";

export const ActivitiesList: React.FC<IResourceComponentsProps> = () => {
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
        Activity is a section where you can determine after which rating the
        user will receive a particular SBT. You must select a prototype and
        other parameters, which will then be used for automatic minting. As soon
        as the user reaches a certain rating, he will automatically receive a
        new NFT,
        <span
          style={{
            color: "red",
          }}
        >
          Attention! Token threshold must be greater than 0 and this is
          equivalent for user rating.
        </span>
      </Space>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="Id" />
        <Table.Column dataIndex="name" title="Name" />
        <Table.Column dataIndex="description" title="Description" />
        <Table.Column dataIndex="token_threshold" title="Token Threshold" />
        <Table.Column
          dataIndex="sbt_prototype_id"
          title="Prototype NFT"
          render={(_, record: BaseRecord) => (
            <Link to={`/prototype-nfts/show/${record.sbt_prototype_id}`}>
              Link to Prototype NFT
            </Link>
          )}
        />

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
