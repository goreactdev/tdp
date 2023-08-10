import React from "react";
import {
  IResourceComponentsProps,
  useApiUrl,
  useList,
  useNavigation,
} from "@refinedev/core";
import { getValueFromEvent, useForm } from "@refinedev/antd";
import {
  Form,
  Input,
  DatePicker,
  AutoComplete,
  Space,
  Typography,
  Upload,
  Button,
} from "antd";
import { Create } from "../../components/crud/create";


export const PrototypeNftsCreate: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps, onFinish } = useForm();

  const apiUrl = useApiUrl();

  const [attributes, setAttributes] = React.useState<
    {
      trait_type: string;
      value: string;
    }[]
  >([]);

  
  const handleOnFinish = (values: any) => {
    onFinish({
      ...values,
      attributes,
    });
};
  return (
    <Create  saveButtonProps={saveButtonProps}
    
    text='Create NFT Template'
    >
      <Form  {...formProps} onFinish={handleOnFinish} layout="vertical">
        <Form.Item
          label="Display Name"
          name={["display_name"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input size="large" placeholder="Display Name" />
        </Form.Item>

        <Form.Item
          label="Description"
          name={["description"]}
          rules={[
            {
              required: true,
            },
          ]}
        >
          <Input size="large" placeholder="Description" />
        </Form.Item>

        <Form.Item label="Image">
          <Form.Item
            name="image"
            valuePropName="fileList"
            getValueFromEvent={getValueFromEvent}
            noStyle
            rules={[
              {
                required: true,
              },
            ]}  
          >
            <Upload.Dragger
              name="file"
              headers={{
                Authorization: `Bearer ${localStorage.getItem("refine-auth")}`,
              }}
              action={`${apiUrl}/media/upload`}
              listType="picture"
              maxCount={1}
              multiple
              
            >
              <p className="ant-upload-text">Drag & drop a file in this area</p>
            </Upload.Dragger>
          </Form.Item>
        </Form.Item>
        <Form.Item label="Rating points" name={["weight"]}>
          <Input size="large" placeholder="Rating points of the item" />
        </Form.Item>
        <h1>Attributes</h1>
          {attributes.map((attr, index) => (
            <div key={index}>
              <Form.Item label={`Key ${index + 1}`} 
              
              >
                <Input
                  size="large"
                  placeholder={`Key ${index + 1}`}
                  value={attr.trait_type}
                  onChange={(e) =>
                    setAttributes(
                      attributes.map((item, i) =>
                        i === index ? { ...item, trait_type: e.target.value } : item
                      )
                    )
                  }
                />
              </Form.Item>
              <Form.Item label={`Value ${index + 1}`}
              
              >
                <Input
                  size="large"
                  placeholder={`Value ${index + 1}`}
                  value={attr.value}
                  onChange={(e) =>
                    setAttributes(
                      attributes.map((item, i) =>
                        i === index ? { ...item, value: e.target.value } : item
                      )
                    )
                  }
                />
              </Form.Item>
            </div>
          ))}
        <Space size="large">
          <Button
            type="primary"
            onClick={() =>{
              // add empty attribute
              setAttributes([
                ...attributes,
                {
                  trait_type: "",
                  value: "",
                },
              ])
            }
            }
          >
            Add Attribute
          </Button>
          {attributes.length > 0 && (
            <Button
              
              onClick={() =>
                // remove last attribute
                setAttributes(attributes.slice(0, attributes.length - 1))
              }
            >
              Remove Attribute
            </Button>
          )}
        </Space>
      </Form>
    </Create>
  );
};
