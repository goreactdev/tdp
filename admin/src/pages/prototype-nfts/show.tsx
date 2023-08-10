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

export const PrototypeNftsShow: React.FC<IResourceComponentsProps> = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;

  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>Id</Title>
      <NumberField value={record?.id ?? ""} />
      <Title level={5}>Base64</Title>
      <TextField value={record?.base64} />
      <Title level={5}>Name</Title>
      <TextField value={record?.name} />
      <Title level={5}>Description</Title>
      <TextField value={record?.description} />
      <Title level={5}>External Url</Title>
      <UrlField value={record?.external_url} />
      <Title level={5}>Image</Title>
      <ImageField style={{ maxWidth: 200 }} value={record?.image} />
      <Title level={5}>Marketplace</Title>
      <UrlField value={record?.marketplace} />
      <Title level={5}>Attributes</Title>
      {record?.attributes?.map((item: { trait_type: string; value: string }) => (
        <TagField
          value={item.trait_type + " : " + item.value}
          key={item?.trait_type}
        />
      ))}
    </Show>
  );
};
