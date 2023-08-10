import React from "react";
import { IResourceComponentsProps, useShow } from "@refinedev/core";
import {
  Show,
  NumberField,
  TagField,
  TextField,
  UrlField,
  ImageField,
} from "@refinedev/antd";
import { Typography } from "antd";

const { Title } = Typography;

export const CollectionsShow: React.FC<IResourceComponentsProps> = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;

  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>Id</Title>
      <NumberField value={record?.id ?? ""} />
      <Title level={5}>Raw Address</Title>
      <TextField value={record?.raw_address} />
      <Title level={5}>Friendly Address</Title>
      <TextField value={record?.friendly_address} />
      <Title level={5}>Next Item Index</Title>
      <NumberField value={record?.next_item_index ?? ""} />
      <Title level={5}>Content Uri</Title>
      <UrlField value={record?.content_uri} />
      <Title level={5}>Raw Owner Address</Title>
      <TextField value={record?.raw_owner_address} />
      <Title level={5}>Friendly Owner Address</Title>
      <TextField value={record?.friendly_owner_address} />
      <Title level={5}>Name</Title>
      <TextField value={record?.name} />
      <Title level={5}>Description</Title>
      <TextField value={record?.description} />
      <Title level={5}>Image</Title>
      <ImageField style={{ maxWidth: 200 }} value={record?.image} />
      <Title level={5}>Content Json</Title>
      <TextField value={JSON.stringify(record?.content_json)} />
      <Title level={5}>Default rating points</Title>
      <NumberField value={record?.default_weight ?? ""} />
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
      <Title level={5}>Version</Title>
      <NumberField value={record?.version ?? ""} />
    </Show>
  );
};
