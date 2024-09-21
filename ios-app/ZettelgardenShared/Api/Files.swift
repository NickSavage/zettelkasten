import Foundation
import SwiftUI

public func fetchFile(
    session: HttpSession,
    fileId: Int,
    originalFileName: String,
    completion: @escaping (Result<URL, Error>) -> Void
) {
    guard let url = URL(string: session.environment + "/files/download/" + String(fileId)) else {
        completion(.failure(NetworkError.invalidURL))
        return
    }
    let token = session.token ?? ""

    performFileDownloadRequest(
        with: url,
        token: token,
        originalFileName: originalFileName,
        completion: completion
    )
}

public func fetchFiles(
    session: HttpSession,
    completion: @escaping (Result<[File], Error>) -> Void
) {

    guard let url = URL(string: session.environment + "/files") else {
        completion(.failure(NetworkError.invalidURL))
        return
    }
    let token = session.token ?? ""

    performRequest(with: url, token: token, completion: completion)

}

public func uploadFileImplementation(
    fileURL: URL,
    cardPK: Int,
    session: HttpSession,
    completion: @escaping (Result<UploadFileResponse, Error>) -> Void
) {
    let baseUrl: String = session.environment
    let token = session.token ?? ""

    let urlString = "\(baseUrl)/files/upload"
    guard let url = URL(string: urlString) else {
        completion(.failure(NetworkError.invalidURL))
        return
    }

    var request = URLRequest(url: url)
    request.httpMethod = "POST"
    request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

    // Create the boundary and Content-Type
    let boundary = UUID().uuidString
    request.setValue(
        "multipart/form-data; boundary=\(boundary)",
        forHTTPHeaderField: "Content-Type"
    )

    // Create the form data
    var body = Data()
    if let fileData = try? Data(contentsOf: fileURL) {
        let maxSize = 10 * 1024 * 1024  // 10 MB
        guard fileData.count <= maxSize else {
            completion(.failure(NetworkError.requestFailed))  // Custom error for file size
            return
        }

        // Append file data
        body.append("\r\n--\(boundary)\r\n".data(using: .utf8)!)
        body.append(
            "Content-Disposition: form-data; name=\"file\"; filename=\"\(fileURL.lastPathComponent)\"\r\n"
                .data(using: .utf8)!
        )
        body.append("Content-Type: application/octet-stream\r\n\r\n".data(using: .utf8)!)
        body.append(fileData)

        // Append card_pk
        body.append("\r\n--\(boundary)\r\n".data(using: .utf8)!)
        body.append("Content-Disposition: form-data; name=\"card_pk\"\r\n\r\n".data(using: .utf8)!)
        body.append("\(cardPK)".data(using: .utf8)!)
    }

    body.append("\r\n--\(boundary)--\r\n".data(using: .utf8)!)
    request.httpBody = body

    // Perform the request
    let task = URLSession.shared.dataTask(with: request) { data, response, error in
        if let error = error {
            completion(.failure(error))
            return
        }

        guard let httpResponse = response as? HTTPURLResponse else {
            completion(.failure(NetworkError.requestFailed))
            return
        }

        // Check for a successful HTTP status code.
        if (200...299).contains(httpResponse.statusCode) {
            // Expect data when the request succeeds
            guard let data = data else {
                completion(.failure(NetworkError.requestFailed))
                return
            }

            do {
                let uploadResponse = try JSONDecoder().decode(UploadFileResponse.self, from: data)
                completion(.success(uploadResponse))
            }
            catch {
                completion(.failure(NetworkError.decodingError(error)))
            }
        }
        else {
            // Attempt to parse server-provided error information
            if let data = data,
                let serverError = try? JSONSerialization.jsonObject(with: data, options: [])
                    as? [String: Any]
            {
                print("Server error:", serverError)
                // Capture detailed error information for debugging/logging
                let errorMessage = serverError["message"] as? String ?? "Unknown server error"
                completion(
                    .failure(
                        NSError(
                            domain: "ServerError",
                            code: httpResponse.statusCode,
                            userInfo: [NSLocalizedDescriptionKey: errorMessage]
                        )
                    )
                )
            }
            else {
                // Fallback error message if parsing fails
                completion(
                    .failure(
                        NSError(
                            domain: "ServerError",
                            code: httpResponse.statusCode,
                            userInfo: [
                                NSLocalizedDescriptionKey:
                                    "Unknown error with code \(httpResponse.statusCode)"
                            ]
                        )
                    )
                )
            }
        }
    }

    task.resume()
}
