import React from "react";
import { IResourceComponentsProps, useShow, useOne } from "@refinedev/core";
import { Show, NumberField, TextField } from "@refinedev/antd";
import { Space, Typography } from "antd";
import { Link } from "react-router-dom";

const { Title } = Typography;

export const RewardsShow: React.FC<IResourceComponentsProps> = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;

  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Space direction="vertical" size={"middle"}>
        <div>
          <Title level={5}>Id</Title>
          <NumberField value={record?.id ?? ""} />
        </div>
        <div>
          <Title level={5}>SBT</Title>
          <Link to={`/minted-nfts/show/${record?.sbt_token_id}`}>
            Link to SBT
          </Link>
        </div>
        <div>
          <Title level={5}>User</Title>
          <Link to={`/users/show/${record?.user_id}`}>Link to User</Link>
        </div>

        <div>
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
        </div>
        <div>
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
        </div>

        <div>
          <Title level={5}>Version</Title>
          <NumberField value={record?.version ?? ""} />
        </div>
      </Space>
    </Show>
  );
};
