import React from "react";
import { IResourceComponentsProps, BaseRecord } from "@refinedev/core";
import {
  useTable,
  List,
  EditButton,
  ShowButton,
  TagField,
  TextField,
  ImageField,
} from "@refinedev/antd";
import { Table, Space } from "antd";
import { CreateButton } from "../../components/buttons/create";
import { LinkedAccount } from "../../authProvider";

export const UsersList: React.FC<IResourceComponentsProps> = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="Id" />
        <Table.Column dataIndex="first_name" title="First Name" />
        <Table.Column dataIndex="last_name" title="Last Name" />
        <Table.Column dataIndex="username" title="Username" />
        <Table.Column
          dataIndex="role"
          title="Role"
          render={(value: BaseRecord) => <>{value?.name || "No Role"}</>}
        />

        <Table.Column dataIndex="awards_count" title="Awards Count" />
        <Table.Column
          dataIndex="avatar_url"
          render={(value) => {
            return <ImageField value={value} style={{ maxWidth: 100 }} />;
          }}
          title="Image"
        ></Table.Column>

        <Table.Column
          dataIndex="linked_accounts"
          title="Linked Accounts"
          render={(value: LinkedAccount[]) => (
            <>
              {value?.map((item) => (
                <a
                  href={
                    item.provider === "github"
                      ? "https://github.com/" + item.login
                      : "https://t.me/" + item.login
                  }
                  target="_blank"
                  rel="noreferrer"
                >
                  <TagField
                    value={
                      item.provider === "github"
                        ? "https://github.com/" + item.login
                        : "https://t.me/" + item.login
                    }
                    key={item?.provider}
                  />
                </a>
              ))}
            </>
          )}
        />
        <Table.Column dataIndex="version" title="Version" />
        <Table.Column
          dataIndex="created_at"
          title="Created At"
          render={(value) => (
            <>{new Intl.DateTimeFormat("gb-GB").format(value.created_at)}</>
          )}
        />
        <Table.Column
          dataIndex="updated_at"
          title="Updated At"
          render={(value) => (
            <>{new Intl.DateTimeFormat("gb-GB").format(value.updated_at)}</>
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
