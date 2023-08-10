import React from "react";
import { IResourceComponentsProps, BaseRecord } from "@refinedev/core";
import {
  useTable,
  List,
  EditButton,
  ShowButton,
  ImageField,
  TagField,
} from "@refinedev/antd";
import { Table, Space } from "antd";

export const PrototypeNftsList: React.FC<IResourceComponentsProps> = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  return (
    <List headerProps={
      {
        title: "NFT Templates"
      }
    }>
      <Space style={
        {
            marginBottom: "20px",
            maxWidth: "25rem",
            backgroundColor: "#fff",
            borderRadius: "10px",
            padding: "10px",
            fontWeight: "600",
            fontSize: "14px",
        }}
        direction="vertical"
        size="small">
        NFT Templates will be used to mint NFTs later.
        <span
          style={{
            color: "red",
          }}
        >
        To Mint an NFT open "NFTs" and click "Mint".
        </span>
      </Space>

      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="Id" />
        <Table.Column dataIndex="display_name" title="Display Name" />
        <Table.Column dataIndex="description" title="Description" />
        <Table.Column
          dataIndex={["image"]}
          title="Image"
          render={(value: any) => (
            <ImageField style={{ maxWidth: "100px" }} value={value} />
          )}
        />
        <Table.Column dataIndex="weight" title="Rating Points" />

        <Table.Column
          dataIndex="attributes"
          title="Attributes"
          render={(value: any[]) => (
            <>
              {value?.map((item) => (
                <TagField
                  value={item.trait_type + " : " + item.value}
                  key={item.trait_type}
                />
              ))}
            </>
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
