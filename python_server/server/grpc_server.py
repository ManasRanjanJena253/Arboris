import grpc
import os
import sys
from dotenv import load_dotenv

current_dir = os.path.dirname(os.path.abspath(__file__))
project_root = os.path.abspath(os.path.join(current_dir, "../../"))
generated_python_dir = os.path.join(project_root, "generated", "python")

sys.path.insert(0, generated_python_dir)

from generated.python import arboris_pb2
from generated.python import arboris_pb2_grpc
from concurrent import futures

load_dotenv()



class ReviewServiceServicer(arboris_pb2_grpc.ReviewServiceServicer):
    def __init__(self):
        pass

    def ClassifyDiff(self, request, context):
        return arboris_pb2.ClassifyDiffResponse(
            severity = None,
            confidence = 0.0
        )

    def DetectPatterns(self, request, context):
        return arboris_pb2.DetectPatternsResponse(
            matched = False,
            pattern_id = "23455432",
            confidence = 0.0
        )

    def GenerateEmbedding(self, request, context):
        return arboris_pb2.GenerateEmbeddingResponse(
            embedding = arboris_pb2.EmbeddingVector(
                values = [0.0, 1.0],
                model_name = "mock_model"
            ),
        )

    def ReviewPullRequest(self, request, context):
        return None

def serve():
    PORT = os.getenv("PYTHON_SERVER_PORT")
    HOST = os.getenv("PYTHON_SERVER_HOST")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers = 10))

    arboris_pb2_grpc.add_ReviewServiceServicer_to_server(ReviewServiceServicer(), server)
    server.add_insecure_port(f"{HOST}:{PORT}")
    print(f"Starting the server at port {HOST}:{PORT}")

    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
