import React from "react";
import { IResourceComponentsProps, BaseRecord, } from "@refinedev/core";
import { useTable, List, ShowButton, } from "@refinedev/antd";
import { Table, Space, Typography } from "antd";
import { Link } from "react-router-dom";

export const RewardsList: React.FC<IResourceComponentsProps> = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  return (
    <List canCreate={false}>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex={["id"]} title="ID" render={(value) => value} />

        <Table.Column
          dataIndex={["user_id"]}
          title="User"
          render={(value) => (
            <Space direction="vertical">
              <Link to={`/users/show/${value}`}>Click to view</Link>
              <Typography>User ID: {value}</Typography>
            </Space>
          )}
        />
        <Table.Column
          dataIndex={["sbt_token_id"]}
          title="SBT Token"
          render={(value) => (
            <Link to={`/minted-nfts/show/${value}`}>Click to view</Link>
          )}
        />

        <Table.Column dataIndex="weight" title="Points Earned" />
        <Table.Column
          dataIndex="created_at"
          title="Created At"
          render={(value) => new Date(value * 1000).toLocaleDateString("ru-RU")}
        />

        <Table.Column
          title="Actions"
          dataIndex="actions"
          render={(_, record: BaseRecord) => (
            <Space>
              <ShowButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};
