import { useState } from "react";
import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from "@tanstack/react-query";
import {
  Theme,
  Spinner,
  Callout,
  Table,
  Card,
  Heading,
} from "@radix-ui/themes";
import { InfoCircledIcon } from "@radix-ui/react-icons";
import { ErrorBoundary } from "react-error-boundary";
import { Line, LineChart, XAxis, YAxis } from "recharts";

import "@radix-ui/themes/styles.css";
import "./app.css";

// Create a client
const queryClient = new QueryClient();

function ErrorFallback({ error }: { error: Error }) {
  return (
    <div>
      Oh no! Something went wrong:
      <pre>{error.message}</pre>
    </div>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Theme>
        <Page title="Fingerprints">
          <Card className="card table">
            <Heading size="6">Fingerprints</Heading>
            <ErrorBoundary FallbackComponent={ErrorFallback}>
              <Fingerprints />
            </ErrorBoundary>
          </Card>
          <Card className="card">
            <Heading size="6">Amount per 30 seconds</Heading>
            <FingerprintGraphs />
          </Card>
        </Page>
      </Theme>
    </QueryClientProvider>
  );
}

const BASE_URL = "http://localhost:8080/api";

const Page = ({
  children,
  title,
}: {
  children: React.ReactNode;
  title: string;
}) => (
  <div className="page">
    <div className="header">
      <Heading size="8">{title}</Heading>
    </div>
    <div className="content">{children}</div>
  </div>
);

// Data is of the format:
// {
//   "fingerprints": [{ input: string, timestamp: string }]
//   "current_page": int,
//   "total_pages": int,
// }
const fetchFingerprints = async () => {
  const response = await fetch(`${BASE_URL}/fingerprints`);
  return response.json();
};

const Fingerprints = () => {
  const [explode, setExplode] = useState(false);

  const { data, isLoading, error } = useQuery({
    queryKey: ["fingerprints"],
    queryFn: fetchFingerprints,
  });

  if (isLoading) return <Spinner />;
  if (error)
    return (
      <Callout.Root>
        <Callout.Icon>
          <InfoCircledIcon />
        </Callout.Icon>
        <Callout.Text>
          An error occurred while fetching fingerprints. Please try again.
        </Callout.Text>
      </Callout.Root>
    );

  if (explode) {
    throw new Error("A bug was left in during one of 11 backports!");
  }

  return (
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.ColumnHeaderCell>Input</Table.ColumnHeaderCell>
          <Table.ColumnHeaderCell>Timestamp</Table.ColumnHeaderCell>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {data.fingerprints.map(
          (fingerprint: { input: string; timestamp: string }, i: number) => (
            <Table.Row key={i}>
              <Table.Cell className="explode" onClick={() => setExplode(true)}>
                {fingerprint.input}
              </Table.Cell>
              <Table.Cell>{fingerprint.timestamp}</Table.Cell>
            </Table.Row>
          ),
        )}
      </Table.Body>
    </Table.Root>
  );
};

// Data is of the format:
// [{ timestamp: string, count: int }]
const fetchFingerprintCount = async () => {
  const response = await fetch(`${BASE_URL}/fingerprints/count`);
  return response.json();
};

const FingerprintGraphs = () => {
  const { data, isLoading, error } = useQuery({
    queryKey: ["fingerprintCount"],
    queryFn: fetchFingerprintCount,
  });

  if (isLoading) return <Spinner />;
  if (error)
    return (
      <Callout.Root>
        <Callout.Icon>
          <InfoCircledIcon />
        </Callout.Icon>
        <Callout.Text>
          An error occurred while fetching fingerprint count. Please try again.
        </Callout.Text>
      </Callout.Root>
    );

  // Render a rechart line graph.
  return (
    <div>
      <LineChart width={600} height={300} data={data}>
        <XAxis dataKey="timestamp" />
        <YAxis />
        <Line type="monotone" dataKey="count" stroke="#8884d8" />
      </LineChart>
    </div>
  );
};

export default App;
