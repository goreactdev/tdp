import React from "react";
import { IResourceComponentsProps, useShow } from "@refinedev/core";
import {
  Show,
  NumberField,
  TagField,
  TextField,
  ImageField,
} from "@refinedev/antd";
import { Typography } from "antd";
import { LinkedAccount } from "../../authProvider";

const { Title } = Typography;

export const UsersShow: React.FC<IResourceComponentsProps> = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;

  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>Id</Title>
      <NumberField value={record?.id ?? ""} />
      <Title level={5}>First Name</Title>
      <TextField value={record?.first_name} />
      <Title level={5}>Last Name</Title>
      <TextField value={record?.last_name} />
      <Title level={5}>Username</Title>
      <TextField value={record?.username} />
      <Title level={5}>Role</Title>
      <TextField value={record?.role?.name ?? ""} />

      <Title level={5}>Friendly Address</Title>
      <TextField value={record?.friendly_address} />
      <Title level={5}>Raw Address</Title>
      <TextField value={record?.raw_address} />
      <Title level={5}>Job</Title>
      <TextField value={record?.job} />
      <Title level={5}>Bio</Title>
      <TextField value={record?.bio} />
      <Title level={5}>Avatar Url</Title>
      <ImageField style={{ maxWidth: 200 }} value={record?.avatar_url} />
      <Title level={5}>Created At</Title>
      <TextField
        value={new Date(record?.created_at * 1000 ?? "").toLocaleString(
          "ru-RU",
          {
            year: "numeric",
            month: "numeric",
            day: "numeric",
          }
        )}
      />
      <Title level={5}>Updated At</Title>
      <TextField
        value={new Date(record?.updated_at * 1000 ?? "").toLocaleString(
          "ru-RU",
          {
            year: "numeric",
            month: "numeric",
            day: "numeric",
          }
        )}
      />
      <Title level={5}>Awards Count</Title>
      <NumberField value={record?.awards_count ?? ""} />
      <Title level={5}>Linked Accounts</Title>
      {record?.linked_accounts?.map((item: LinkedAccount) => (
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
      <Title level={5}>Version</Title>
      <NumberField value={record?.version ?? ""} />
    </Show>
  );
};
